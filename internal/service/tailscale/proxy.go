package tailscale

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	tcpProtocol = "tcp"
	udpProtocol = "udp"

	udpSessionIdleTimeout = time.Minute
	udpSessionSweepPeriod = 15 * time.Second
	maxUDPSessions        = 256
)

var errUDPSessionLimit = errors.New("UDP session limit reached")

type portForward struct {
	port     uint16
	protocol string
}

type tcpForward struct {
	listenAddr string
	targetAddr string
}

type udpForward struct {
	port       uint16
	targetAddr string
}

type proxyForwards struct {
	tcp []tcpForward
	udp []udpForward
}

type activeTCPForward struct {
	listener net.Listener
	target   string
}

type activeUDPForward struct {
	packetConn net.PacketConn
	target     string
}

type udpSession struct {
	conn       net.Conn
	clientAddr net.Addr
	lastActive atomic.Int64
}

type udpSessionManager struct {
	ctx        context.Context
	packetConn net.PacketConn
	targetAddr string
	log        serviceLogger

	mu       sync.Mutex
	sessions map[string]*udpSession
	wg       sync.WaitGroup
}

func parsePortForwardSpec(value string) (portForward, error) {
	spec := strings.TrimSpace(value)
	parts := strings.Split(spec, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return portForward{}, fmt.Errorf("invalid forward port %q: expected port/protocol", value)
	}

	for _, char := range parts[0] {
		if char < '0' || char > '9' {
			return portForward{}, fmt.Errorf("invalid forward port %q: port must be a decimal number", value)
		}
	}
	port, err := strconv.ParseUint(parts[0], 10, 16)
	if err != nil || port == 0 {
		return portForward{}, fmt.Errorf("invalid forward port %q: port must be between 1 and 65535", value)
	}

	protocol := strings.ToLower(parts[1])
	if protocol != tcpProtocol && protocol != udpProtocol {
		return portForward{}, fmt.Errorf("invalid forward port %q: protocol must be tcp or udp", value)
	}

	return portForward{port: uint16(port), protocol: protocol}, nil
}

func parsePortForwardSpecs(values []string) ([]portForward, error) {
	seen := make(map[portForward]struct{}, len(values))
	forwards := make([]portForward, 0, len(values))
	for _, value := range values {
		forward, err := parsePortForwardSpec(value)
		if err != nil {
			return nil, err
		}
		if _, exists := seen[forward]; exists {
			continue
		}
		seen[forward] = struct{}{}
		forwards = append(forwards, forward)
	}
	slices.SortFunc(forwards, func(a, b portForward) int {
		if a.port != b.port {
			return int(a.port) - int(b.port)
		}
		return strings.Compare(a.protocol, b.protocol)
	})
	return forwards, nil
}

func normalizePortForwardSpecs(values []string) ([]string, error) {
	forwards, err := parsePortForwardSpecs(values)
	if err != nil {
		return nil, err
	}
	normalized := make([]string, 0, len(forwards))
	for _, forward := range forwards {
		normalized = append(normalized, strconv.Itoa(int(forward.port))+"/"+forward.protocol)
	}
	return normalized, nil
}

func buildForwards(localServerAddr string, additionalPorts []string) (proxyForwards, error) {
	httpListenAddr, httpTargetAddr, err := proxyAddresses(localServerAddr)
	if err != nil {
		return proxyForwards{}, err
	}

	configured, err := parsePortForwardSpecs(additionalPorts)
	if err != nil {
		return proxyForwards{}, err
	}

	forwards := proxyForwards{
		tcp: []tcpForward{{listenAddr: httpListenAddr, targetAddr: httpTargetAddr}},
	}
	seenTCP := map[string]struct{}{httpListenAddr: {}}
	for _, forward := range configured {
		port := strconv.Itoa(int(forward.port))
		targetAddr := net.JoinHostPort("localhost", port)
		switch forward.protocol {
		case tcpProtocol:
			listenAddr := net.JoinHostPort("", port)
			if _, exists := seenTCP[listenAddr]; exists {
				continue
			}
			seenTCP[listenAddr] = struct{}{}
			forwards.tcp = append(forwards.tcp, tcpForward{
				listenAddr: listenAddr,
				targetAddr: targetAddr,
			})
		case udpProtocol:
			forwards.udp = append(forwards.udp, udpForward{
				port:       forward.port,
				targetAddr: targetAddr,
			})
		}
	}
	return forwards, nil
}

func proxyAddresses(localServerAddr string) (tailnetListenAddr, localTargetAddr string, err error) {
	host, port, err := net.SplitHostPort(localServerAddr)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse local HTTP server address: %w", err)
	}

	portNumber, err := strconv.ParseUint(port, 10, 16)
	if err != nil || portNumber == 0 {
		return "", "", fmt.Errorf("invalid local HTTP server port %q", port)
	}

	if host == "" {
		host = "127.0.0.1"
	} else if ip := net.ParseIP(host); ip != nil && ip.IsUnspecified() {
		if ip.To4() == nil {
			host = "::1"
		} else {
			host = "127.0.0.1"
		}
	}

	return net.JoinHostPort("", port), net.JoinHostPort(host, port), nil
}

func serveForwards(ctx context.Context, node embeddedNode, forwards proxyForwards, log serviceLogger) error {
	proxyCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	workerCount := 0
	results := make(chan error, 2)
	if len(forwards.tcp) > 0 {
		workerCount++
		go func() {
			results <- serveTCPForwards(proxyCtx, node, forwards.tcp, log)
		}()
	}
	if len(forwards.udp) > 0 {
		workerCount++
		go func() {
			results <- serveUDPForwards(proxyCtx, node, forwards.udp, log)
		}()
	}
	if workerCount == 0 {
		<-ctx.Done()
		return nil
	}

	result := <-results
	cancel()
	for range workerCount - 1 {
		if err := <-results; result == nil {
			result = err
		}
	}
	if ctx.Err() != nil {
		return nil
	}
	return result
}

func serveTCPForwards(ctx context.Context, node embeddedNode, forwards []tcpForward, log serviceLogger) error {
	active := make([]activeTCPForward, 0, len(forwards))
	for _, forward := range forwards {
		listener, err := node.Listen(tcpProtocol, forward.listenAddr)
		if err != nil {
			for _, item := range active {
				_ = item.listener.Close()
			}
			return fmt.Errorf("failed to listen on tailscale address %s: %w", forward.listenAddr, err)
		}
		active = append(active, activeTCPForward{listener: listener, target: forward.targetAddr})
		log.Infof("embedded tailscale node is forwarding TCP %s to %s", listener.Addr(), forward.targetAddr)
	}

	proxyCtx, cancel := context.WithCancel(ctx)
	results := make(chan error, len(active))
	for _, item := range active {
		go func() {
			results <- serveTCPProxy(proxyCtx, item.listener, item.target, log)
		}()
	}

	var result error
	remaining := len(active)
	select {
	case <-ctx.Done():
	case result = <-results:
		remaining--
	}
	cancel()
	for range remaining {
		if err := <-results; result == nil {
			result = err
		}
	}
	if ctx.Err() != nil {
		return nil
	}
	return result
}

func serveTCPProxy(ctx context.Context, listener net.Listener, targetAddr string, log serviceLogger) error {
	proxyCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	stopClose := context.AfterFunc(proxyCtx, func() {
		_ = listener.Close()
	})
	defer stopClose()

	var connections sync.WaitGroup
	defer func() {
		cancel()
		connections.Wait()
	}()

	for {
		incoming, err := listener.Accept()
		if err != nil {
			if proxyCtx.Err() != nil {
				return nil
			}
			return fmt.Errorf("failed to accept tailscale connection: %w", err)
		}

		connections.Add(1)
		go func() {
			defer connections.Done()
			proxyTCPConnection(proxyCtx, incoming, targetAddr, log)
		}()
	}
}

func proxyTCPConnection(ctx context.Context, incoming net.Conn, targetAddr string, log serviceLogger) {
	defer incoming.Close()

	outgoing, err := (&net.Dialer{}).DialContext(ctx, tcpProtocol, targetAddr)
	if err != nil {
		if ctx.Err() == nil {
			log.Errorf("failed to connect tailscale request to local service: %v", err)
		}
		return
	}
	defer outgoing.Close()

	stopClose := context.AfterFunc(ctx, func() {
		_ = incoming.Close()
		_ = outgoing.Close()
	})
	defer stopClose()

	copyDone := make(chan struct{}, 2)
	go func() {
		_, _ = io.Copy(outgoing, incoming)
		copyDone <- struct{}{}
	}()
	go func() {
		_, _ = io.Copy(incoming, outgoing)
		copyDone <- struct{}{}
	}()
	<-copyDone
}

func serveUDPForwards(ctx context.Context, node embeddedNode, forwards []udpForward, log serviceLogger) error {
	ipv4, ipv6 := node.TailscaleIPs()
	ips := make([]string, 0, 2)
	if ipv4.IsValid() {
		ips = append(ips, ipv4.String())
	}
	if ipv6.IsValid() {
		ips = append(ips, ipv6.String())
	}
	if len(ips) == 0 {
		return errors.New("embedded tailscale node has no assigned IP address for UDP forwarding")
	}

	active := make([]activeUDPForward, 0, len(forwards)*len(ips))
	for _, forward := range forwards {
		port := strconv.Itoa(int(forward.port))
		for _, ip := range ips {
			listenAddr := net.JoinHostPort(ip, port)
			packetConn, err := node.ListenPacket(udpProtocol, listenAddr)
			if err != nil {
				for _, item := range active {
					_ = item.packetConn.Close()
				}
				return fmt.Errorf("failed to listen on tailscale UDP address %s: %w", listenAddr, err)
			}
			active = append(active, activeUDPForward{packetConn: packetConn, target: forward.targetAddr})
			log.Infof("embedded tailscale node is forwarding UDP %s to %s", packetConn.LocalAddr(), forward.targetAddr)
		}
	}

	proxyCtx, cancel := context.WithCancel(ctx)
	results := make(chan error, len(active))
	for _, item := range active {
		go func() {
			results <- serveUDPProxy(proxyCtx, item.packetConn, item.target, log)
		}()
	}

	var result error
	remaining := len(active)
	select {
	case <-ctx.Done():
	case result = <-results:
		remaining--
	}
	cancel()
	for range remaining {
		if err := <-results; result == nil {
			result = err
		}
	}
	if ctx.Err() != nil {
		return nil
	}
	return result
}

func serveUDPProxy(ctx context.Context, packetConn net.PacketConn, targetAddr string, log serviceLogger) error {
	proxyCtx, cancel := context.WithCancel(ctx)
	manager := &udpSessionManager{
		ctx:        proxyCtx,
		packetConn: packetConn,
		targetAddr: targetAddr,
		log:        log,
		sessions:   make(map[string]*udpSession),
	}

	stopClose := context.AfterFunc(proxyCtx, func() {
		_ = packetConn.Close()
	})
	sweepDone := make(chan struct{})
	go func() {
		defer close(sweepDone)
		manager.sweepIdleSessions()
	}()
	defer func() {
		cancel()
		stopClose()
		_ = packetConn.Close()
		manager.close()
		<-sweepDone
	}()

	buffer := make([]byte, 64*1024)
	for {
		n, clientAddr, err := packetConn.ReadFrom(buffer)
		if err != nil {
			if proxyCtx.Err() != nil {
				return nil
			}
			return fmt.Errorf("failed to read tailscale UDP packet: %w", err)
		}

		session, err := manager.getSession(clientAddr)
		if errors.Is(err, errUDPSessionLimit) {
			continue
		}
		if err != nil {
			if proxyCtx.Err() == nil {
				log.Errorf("failed to connect tailscale UDP request to local service: %v", err)
			}
			continue
		}

		session.lastActive.Store(time.Now().UnixNano())
		if _, err := session.conn.Write(buffer[:n]); err != nil {
			manager.removeSession(udpSessionKey(clientAddr), session)
			if proxyCtx.Err() == nil {
				log.Errorf("failed to forward tailscale UDP request to local service: %v", err)
			}
		}
	}
}

func (m *udpSessionManager) getSession(clientAddr net.Addr) (*udpSession, error) {
	key := udpSessionKey(clientAddr)
	m.mu.Lock()
	if session := m.sessions[key]; session != nil {
		m.mu.Unlock()
		return session, nil
	}
	if len(m.sessions) >= maxUDPSessions {
		m.mu.Unlock()
		return nil, errUDPSessionLimit
	}
	m.mu.Unlock()

	conn, err := (&net.Dialer{}).DialContext(m.ctx, udpProtocol, m.targetAddr)
	if err != nil {
		return nil, err
	}
	session := &udpSession{conn: conn, clientAddr: clientAddr}
	session.lastActive.Store(time.Now().UnixNano())

	m.mu.Lock()
	if existing := m.sessions[key]; existing != nil {
		m.mu.Unlock()
		_ = conn.Close()
		return existing, nil
	}
	if len(m.sessions) >= maxUDPSessions {
		m.mu.Unlock()
		_ = conn.Close()
		return nil, errUDPSessionLimit
	}
	m.sessions[key] = session
	m.wg.Add(1)
	m.mu.Unlock()

	go m.forwardResponses(key, session)
	return session, nil
}

func (m *udpSessionManager) forwardResponses(key string, session *udpSession) {
	defer m.wg.Done()
	defer m.removeSession(key, session)

	buffer := make([]byte, 64*1024)
	for {
		n, err := session.conn.Read(buffer)
		if err != nil {
			if m.ctx.Err() == nil {
				m.log.Errorf("failed to read local UDP response: %v", err)
			}
			return
		}
		session.lastActive.Store(time.Now().UnixNano())
		if _, err := m.packetConn.WriteTo(buffer[:n], session.clientAddr); err != nil {
			if m.ctx.Err() == nil {
				m.log.Errorf("failed to return UDP response over tailscale: %v", err)
			}
			return
		}
	}
}

func (m *udpSessionManager) sweepIdleSessions() {
	ticker := time.NewTicker(udpSessionSweepPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-m.ctx.Done():
			return
		case now := <-ticker.C:
			cutoff := now.Add(-udpSessionIdleTimeout).UnixNano()
			m.mu.Lock()
			var expired []*udpSession
			for key, session := range m.sessions {
				if session.lastActive.Load() < cutoff {
					delete(m.sessions, key)
					expired = append(expired, session)
				}
			}
			m.mu.Unlock()
			for _, session := range expired {
				_ = session.conn.Close()
			}
		}
	}
}

func (m *udpSessionManager) removeSession(key string, expected *udpSession) {
	m.mu.Lock()
	session := m.sessions[key]
	if session == expected {
		delete(m.sessions, key)
	}
	m.mu.Unlock()
	if session == expected {
		_ = session.conn.Close()
	}
}

func (m *udpSessionManager) close() {
	m.mu.Lock()
	sessions := make([]*udpSession, 0, len(m.sessions))
	for key, session := range m.sessions {
		delete(m.sessions, key)
		sessions = append(sessions, session)
	}
	m.mu.Unlock()
	for _, session := range sessions {
		_ = session.conn.Close()
	}
	m.wg.Wait()
}

func udpSessionKey(addr net.Addr) string {
	return addr.Network() + "\x00" + addr.String()
}
