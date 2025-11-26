package forwarder

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime/debug"
	"strings"
	"time"
	"unsafe"

	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/pkg/logger"
)

func (s *ForwarderServiceImpl) handleInterrupt() {
	s.wg.Done()
}

func (s *ForwarderServiceImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(context.Background())
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.listenHost, s.listenPort))
	if err != nil {
		return fmt.Errorf("failed to listen on %s:%d: %w", s.listenHost, s.listenPort, err)
	}

	s.listener = listener
	logger.GetLogger(ID).Infof("service forwarder is listening on %s:%d", s.listenHost, s.listenPort)

	go func() {
		s.status.SetStartedAt(s.timeSource.Now())
		s.status.SetIsRunning(true)
		defer func() {
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectly crashed, recovered from panic: %v\n%s", r, debug.Stack())
				s.handleInterrupt()
				_ = s.Stop()
			}
		}()

		err = s.hardwareDev.SubscribeRealtime(ID, func(t time.Time, di *explorer.DeviceConfig, dv *explorer.DeviceVariable, cd []explorer.ChannelData) {
			s.messageBusRealtime.Publish(t, di, dv, cd)
		})
		if err != nil {
			logger.GetLogger(ID).Errorf("failed to subscribe to realtime hardware message bus: %v", err)
			return
		}
		err = s.hardwareDev.Subscribe(ID, func(t time.Time, di *explorer.DeviceConfig, dv *explorer.DeviceVariable, cd []explorer.ChannelData) {
			s.messageBus.Publish(t, di, dv, cd)
		})
		if err != nil {
			logger.GetLogger(ID).Errorf("failed to subscribe to hardware message bus: %v", err)
			return
		}

		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					if errors.Is(err, net.ErrClosed) {
						return
					}
					continue
				}
				go s.handleConnection(conn)
			}
		}()

		<-s.ctx.Done()
		s.handleInterrupt()
	}()

	s.wg.Add(1)
	return nil
}

func (a *ForwarderServiceImpl) getChecksum(arr []int32) (checksum uint8) {
	for _, data := range arr {
		bytes := (*[4]byte)(unsafe.Pointer(&data))[:]
		for j := 0; j < int(unsafe.Sizeof(int32(0))); j++ {
			checksum ^= bytes[j]
		}
	}

	return checksum
}

func (a *ForwarderServiceImpl) getDataBytes(tm time.Time, sampleRate int, channelData []explorer.ChannelData) []byte {
	var dataBytes []byte
	for _, channel := range channelData {
		dataStr := strings.Trim(strings.ReplaceAll(fmt.Sprint(channel.Data), " ", ","), "[]")
		msg := fmt.Sprintf(
			"$%d,%s,%s,%s,%s,%d,%d,%s,*%02X\r\n",
			channel.ChannelId,
			a.networkCode,
			a.stationCode,
			a.locationCode,
			channel.ChannelCode,
			tm.UnixMilli(),
			sampleRate,
			dataStr,
			a.getChecksum(channel.Data),
		)

		dataBytes = append(dataBytes, []byte(msg)...)
	}
	return dataBytes
}

func (a *ForwarderServiceImpl) handleConnection(conn net.Conn) {
	key := conn.RemoteAddr().String()
	defer conn.Close()

	logger.GetLogger(ID).Infof("%s - client connected to forwarder service", key)
	defer logger.GetLogger(ID).Infof("%s - client disconnected from forwarder service", key)

	subscribeNormal := func() {
		a.messageBus.Subscribe(key, func(t time.Time, di *explorer.DeviceConfig, dv *explorer.DeviceVariable, cd []explorer.ChannelData) {
			_, err := conn.Write(a.getDataBytes(t, di.GetSampleRate(), cd))
			if err != nil {
				logger.GetLogger(ID).Errorln(err)
			}
		})
	}
	subscribeRealtime := func() {
		a.messageBusRealtime.Subscribe(key, func(t time.Time, di *explorer.DeviceConfig, dv *explorer.DeviceVariable, cd []explorer.ChannelData) {
			_, err := conn.Write(a.getDataBytes(t, di.GetSampleRate(), cd))
			if err != nil {
				logger.GetLogger(ID).Errorln(err)
			}
		})
	}
	unsubscribeAll := func() {
		a.messageBus.Unsubscribe(key)
		a.messageBusRealtime.Unsubscribe(key)
	}

	subscribeNormal()
	buf := make([]byte, 64)
	for useRealtime := false; ; {
		n, err := conn.Read(buf)
		if err != nil {
			unsubscribeAll()
			return
		}

		cmd := string(buf[:n])
		switch cmd {
		case "AT+REALTIME=1\r\n":
			if !useRealtime {
				unsubscribeAll()
				subscribeRealtime()
				useRealtime = true
				logger.GetLogger(ID).Infof("%s switched to REALTIME bus", key)
			}
		case "AT+REALTIME=0\r\n":
			if useRealtime {
				unsubscribeAll()
				subscribeNormal()
				useRealtime = false
				logger.GetLogger(ID).Infof("%s switched to NORMAL bus", key)
			}
		}
	}
}
