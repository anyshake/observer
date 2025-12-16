package ntp_server

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"net/netip"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/anyshake/observer/pkg/logger"
)

func (s *NtpServerServiceImpl) handleInterrupt() {
	s.wg.Done()
}

func (s *NtpServerServiceImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(context.Background())
	}

	ntpSrv, err := newNtpServer(s.listenHost, s.listenPort, func() time.Time {
		return s.timeSource.Now().Add(time.Duration(s.delayControl) * time.Microsecond)
	})
	if err != nil {
		return fmt.Errorf("failed to start NTP server: %w", err)
	}

	go func() {
		s.status.SetStartedAt(s.timeSource.Now())
		s.status.SetIsRunning(true)
		defer func() {
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectly stopped, recovered from panic: %v\n%s", r, debug.Stack())
				s.handleInterrupt()
				_ = s.Stop()
			}
		}()

		logger.GetLogger(ID).Infof("starting NTP server service, listening on %s:%d", s.listenHost, s.listenPort)
		if err := ntpSrv.Run(s.ctx); err != nil {
			logger.GetLogger(ID).Errorf("failed to run NTP server service: %v", err)
			s.handleInterrupt()
			_ = s.Stop()
		}
		s.wg.Done()
	}()

	s.wg.Add(1)
	return nil
}

const (
	LI_NO_WARNING      = 0x0
	LI_ALARM_CONDITION = 0x3
)

const (
	VN_FIRST = 0x1
	VN_LAST  = 0x4
)

const (
	MODE_SYMMETRIC_ACTIVE = 0x1
	MODE_CLIENT           = 0x3
)

const (
	REFERENCE_ID      = "GNSS"
	FROM_1900_TO_1970 = 2208988800
)

type ntpServer struct {
	conn      *net.UDPConn
	wg        sync.WaitGroup
	closeOnce sync.Once
	timeFn    func() time.Time
}

func newNtpServer(addr string, port int, timeFn func() time.Time) (*ntpServer, error) {
	listenAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return nil, err
	}

	c, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		return nil, err
	}

	if timeFn == nil {
		timeFn = time.Now
	}

	return &ntpServer{conn: c, timeFn: timeFn}, nil
}

func (p *ntpServer) Run(ctx context.Context) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := p.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond)); err != nil {
		return err
	}

	for readBuf := make([]byte, 512); ; {
		select {
		case <-ctx.Done():
			p.wg.Wait()
			p.closeOnce.Do(func() { p.conn.Close() })
			return nil
		default:
		}

		n, remoteAddr, err := p.conn.ReadFromUDP(readBuf)
		t2 := p.timeFn()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				if err := p.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond)); err != nil {
					return err
				}
				continue
			}
			return err
		}

		if n > 0 {
			p.handleData(t2, remoteAddr, readBuf[:n])
		}
	}
}

func (p *ntpServer) write(data []byte, laddr *net.UDPAddr) error {
	_, err := p.conn.WriteTo(data, laddr)
	if err != nil {
		return err
	}

	return nil
}

func (p *ntpServer) handleData(t2 time.Time, addr net.Addr, data []byte) {
	addrPort, err := netip.ParseAddrPort(addr.String())
	if err != nil {
		logger.GetLogger(ID).Errorf("failed to parse address and port: %v", err)
		return
	}

	addrPortStr := addrPort.String()
	ipAddr := addrPort.Addr()

	if p.isDataValid(ipAddr, data) {
		logger.GetLogger(ID).Infof("accepted connection from %s", addrPortStr)

		resp, err := p.encodePacket(t2, data)
		if err != nil {
			logger.GetLogger(ID).Errorf("failed to encode NTP response packet: %v", err)
			return
		}

		if err := p.write(resp, addr.(*net.UDPAddr)); err != nil {
			logger.GetLogger(ID).Errorf("write NTP response error: %v", err)
		}
	} else {
		logger.GetLogger(ID).Warnf("rejected connection from %s", addrPortStr)
	}
}

func (p *ntpServer) encodePacket(t2 time.Time, req []byte) ([]byte, error) {
	secT2 := uint32(t2.Unix() + FROM_1900_TO_1970)
	// convert nanoseconds to 32-bit fraction of a second
	fracT2 := uint32((uint64(t2.Nanosecond()) * (1 << 32)) / 1_000_000_000)

	res := make([]byte, 48)
	vn := req[0] & 0x38
	res[0] = vn | 0x04 // version + mode (server)
	res[1] = 1         // stratum
	res[2] = req[2]    // poll
	res[3] = 0xEC      // precision

	// 4-byte reference ID
	copy(res[12:16], []byte(REFERENCE_ID[:4]))

	// reference timestamp (seconds)
	binary.BigEndian.PutUint32(res[16:20], secT2)

	// originate timestamp: client's transmit timestamp (bytes 40..48 of request)
	copy(res[24:32], req[40:48])

	// receive timestamp (T2)
	binary.BigEndian.PutUint32(res[32:36], secT2)
	binary.BigEndian.PutUint32(res[36:40], fracT2)

	// transmit timestamp (T3)
	t3 := p.timeFn()
	secT3 := uint32(t3.Unix() + FROM_1900_TO_1970)
	fracT3 := uint32((uint64(t3.Nanosecond()) * (1 << 32)) / 1_000_000_000)
	binary.BigEndian.PutUint32(res[40:44], secT3)
	binary.BigEndian.PutUint32(res[44:48], fracT3)

	return res, nil
}

func (d *ntpServer) isDataValid(ipAddr netip.Addr, req []byte) bool {
	if len(req) < 48 {
		return false
	}

	li := req[0] >> 6
	ver := (req[0] >> 3) & 0x7
	mode := req[0] & 0x7

	compatibleMode := ipAddr.IsPrivate() || ipAddr.IsLoopback() || ipAddr.IsUnspecified()

	return (li == LI_NO_WARNING || li == LI_ALARM_CONDITION) &&
		(ver >= VN_FIRST && ver <= VN_LAST) &&
		(mode == MODE_CLIENT || (mode == MODE_SYMMETRIC_ACTIVE && compatibleMode))
}
