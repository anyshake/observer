package forwarder

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime/debug"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/ctxutil"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/message"
)

const (
	forwarderWriteTimeout     = 5 * time.Second
	forwarderAcceptRetryDelay = 100 * time.Millisecond
)

type stationCodes struct {
	networkCode  string
	stationCode  string
	locationCode string
}

func (s *ForwarderServiceImpl) Start(ctx context.Context) error {
	if err := s.lifecycle.Lock(ctx); err != nil {
		return err
	}
	defer s.lifecycle.Unlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.GetIsRunning() {
		return errors.New("forwarder service is already running")
	}
	if err := s.parentCtx.Err(); err != nil {
		return err
	}
	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(s.parentCtx)
	}
	s.messageBus = message.NewBus[explorer.Event](ID)
	s.messageBusRealtime = message.NewBus[explorer.Event](ID + "_realtime")
	messageBus := s.messageBus
	messageBusRealtime := s.messageBusRealtime

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.listenHost, s.listenPort))
	if err != nil {
		return fmt.Errorf("failed to listen on %s:%d: %w", s.listenHost, s.listenPort, err)
	}

	s.listener = listener
	logger.GetLogger(ID).Infof("service forwarder is listening on %s:%d", s.listenHost, s.listenPort)

	runCtx, runCancel := s.ctx, s.cancelFn
	codes := stationCodes{networkCode: s.networkCode, stationCode: s.stationCode, locationCode: s.locationCode}
	s.status.SetStartedAt(s.timeSource.Now())
	s.status.SetIsRunning(true)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		var connections sync.WaitGroup
		defer func() {
			runCancel()
			_ = listener.Close()
			connections.Wait()
			messageBus.Close()
			messageBusRealtime.Close()
			_ = s.hardwareDev.UnsubscribeRealtime(ID)
			_ = s.hardwareDev.Unsubscribe(ID)
			s.status.SetStoppedAt(s.timeSource.Now())
			s.status.SetIsRunning(false)
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectedly crashed, recovered from panic: %v\n%s", r, debug.Stack())
			}
		}()

		if err := s.hardwareDev.SubscribeRealtime(ID, func(event explorer.Event) {
			defer service.RecoverCallback(ID, runCancel)
			messageBusRealtime.Publish(event)
		}); err != nil {
			logger.GetLogger(ID).Errorf("failed to subscribe to realtime hardware message bus: %v", err)
			return
		}
		if err := s.hardwareDev.Subscribe(ID, func(event explorer.Event) {
			defer service.RecoverCallback(ID, runCancel)
			messageBus.Publish(event)
		}); err != nil {
			logger.GetLogger(ID).Errorf("failed to subscribe to hardware message bus: %v", err)
			return
		}

		connections.Go(func() {
			defer service.RecoverWorker(ID, runCancel)
			for {
				conn, err := listener.Accept()
				if err != nil {
					if errors.Is(err, net.ErrClosed) || runCtx.Err() != nil {
						return
					}
					logger.GetLogger(ID).Warnf("failed to accept forwarder connection: %v", err)
					select {
					case <-runCtx.Done():
						return
					case <-time.After(forwarderAcceptRetryDelay):
						continue
					}
				}
				connections.Go(func() {
					defer service.RecoverWorker(ID, runCancel)
					s.handleConnection(runCtx, conn, codes, messageBus, messageBusRealtime)
				})
			}
		})

		<-runCtx.Done()
	}()

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

func (a *ForwarderServiceImpl) getDataBytes(codes stationCodes, tm time.Time, sampleRate int, channelData []explorer.ChannelData) []byte {
	var dataBytes []byte
	for _, channel := range channelData {
		dataStr := strings.Trim(strings.ReplaceAll(fmt.Sprint(channel.Data), " ", ","), "[]")
		msg := fmt.Sprintf(
			"$%d,%s,%s,%s,%s,%d,%d,%s,*%02X\r\n",
			channel.ChannelId,
			codes.networkCode,
			codes.stationCode,
			codes.locationCode,
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

func (a *ForwarderServiceImpl) handleConnection(ctx context.Context, conn net.Conn, codes stationCodes, messageBus, messageBusRealtime *message.Bus[explorer.Event]) {
	key := conn.RemoteAddr().String()
	defer conn.Close()
	stopClose := ctxutil.AfterFunc(ctx, func() { _ = conn.Close() })
	defer stopClose()

	logger.GetLogger(ID).Infof("%s - client connected to forwarder service", key)
	defer logger.GetLogger(ID).Infof("%s - client disconnected from forwarder service", key)

	var writeMu sync.Mutex
	writeEvent := func(event explorer.Event) {
		writeMu.Lock()
		defer writeMu.Unlock()

		if err := conn.SetWriteDeadline(time.Now().Add(forwarderWriteTimeout)); err != nil {
			logger.GetLogger(ID).Warnf("%s - failed to set write deadline: %v", key, err)
			_ = conn.Close()
			return
		}
		if _, err := conn.Write(a.getDataBytes(codes, event.Timestamp, event.SampleRate, event.ChannelData)); err != nil {
			logger.GetLogger(ID).Warnf("%s - failed to write forwarded data: %v", key, err)
			_ = conn.Close()
		}
	}
	subscriptionOptions := func(bufferSize int) message.SubscriptionOptions {
		return message.SubscriptionOptions{
			BufferSize: bufferSize,
			Overflow:   message.OverflowDisconnect,
			OnError: func(err error) {
				logger.GetLogger(ID).Warnf("%s - forwarder subscription closed: %v", key, err)
				_ = conn.Close()
			},
		}
	}
	subscribeNormal := func() error {
		return messageBus.Subscribe(key, subscriptionOptions(16), writeEvent)
	}
	subscribeRealtime := func() error {
		return messageBusRealtime.Subscribe(key, subscriptionOptions(32), writeEvent)
	}
	unsubscribeAll := func() {
		_ = messageBus.Unsubscribe(key)
		_ = messageBusRealtime.Unsubscribe(key)
	}
	defer unsubscribeAll()

	if err := subscribeNormal(); err != nil {
		logger.GetLogger(ID).Errorf("%s - failed to subscribe to normal bus: %v", key, err)
		return
	}
	buf := make([]byte, 64)
	for useRealtime := false; ; {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		cmd := string(buf[:n])
		switch cmd {
		case "AT+REALTIME=1\r\n":
			if !useRealtime {
				unsubscribeAll()
				if err := subscribeRealtime(); err != nil {
					logger.GetLogger(ID).Errorf("%s - failed to subscribe to realtime bus: %v", key, err)
					return
				}
				useRealtime = true
				logger.GetLogger(ID).Infof("%s switched to REALTIME bus", key)
			}
		case "AT+REALTIME=0\r\n":
			if useRealtime {
				unsubscribeAll()
				if err := subscribeNormal(); err != nil {
					logger.GetLogger(ID).Errorf("%s - failed to subscribe to normal bus: %v", key, err)
					return
				}
				useRealtime = false
				logger.GetLogger(ID).Infof("%s switched to NORMAL bus", key)
			}
		}
	}
}

func (s *ForwarderServiceImpl) GetListenPort() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.listenPort
}
