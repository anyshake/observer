package frp_client

import (
	"context"
	"fmt"
	"net"
	"runtime/debug"
	"strconv"

	"github.com/anyshake/observer/pkg/logger"
	"github.com/fatedier/frp/client"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/samber/lo"
)

func (s *FrpClientServiceImpl) handleInterrupt() {
	s.wg.Done()
}

func (s *FrpClientServiceImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	localServerHostname, localServerPort, err := s.parseLocalServerAddr()
	if err != nil {
		return err
	}

	falseVal := false
	commonConfig := &v1.ClientCommonConfig{
		LoginFailExit: &falseVal,
		User:          s.user,
		ServerAddr:    s.serverAddr,
		ServerPort:    int(s.serverPort),
		Transport: v1.ClientTransportConfig{
			PoolCount: 5,
			TCPMux:    &s.tcpMux,
			Protocol:  s.transportProtocol,
			TLS: v1.TLSClientConfig{
				Enable:                    &s.tlsEnable,
				DisableCustomTLSFirstByte: &s.disableCustomTLSFirstByte,
			},
		},
		Auth: v1.AuthClientConfig{Token: s.authToken},
	}
	proxyConfig := lo.Ternary[v1.ProxyConfigurer](
		s.useDomainAccess,
		&v1.HTTPProxyConfig{
			ProxyBaseConfig: v1.ProxyBaseConfig{
				Type:         "http",
				Name:         lo.Ternary(s.user != "", fmt.Sprintf("%s.%s", s.user, s.proxyName), s.proxyName),
				ProxyBackend: v1.ProxyBackend{LocalIP: localServerHostname, LocalPort: localServerPort},
				Transport:    v1.ProxyTransport{UseEncryption: s.useEncryption, UseCompression: s.useCompression},
			},
			DomainConfig: v1.DomainConfig{SubDomain: s.subdomain, CustomDomains: s.customDomains},
		},
		&v1.TCPProxyConfig{
			ProxyBaseConfig: v1.ProxyBaseConfig{
				Type:         "tcp",
				Name:         lo.Ternary(s.user != "", fmt.Sprintf("%s.%s", s.user, s.proxyName), s.proxyName),
				ProxyBackend: v1.ProxyBackend{LocalIP: localServerHostname, LocalPort: localServerPort},
				Transport:    v1.ProxyTransport{UseEncryption: s.useEncryption, UseCompression: s.useCompression},
			},
			RemotePort: int(s.remoteOutboundPort),
		},
	)
	svr, err := client.NewService(client.ServiceOptions{
		Common:    commonConfig,
		ProxyCfgs: []v1.ProxyConfigurer{proxyConfig},
	})
	if err != nil {
		return fmt.Errorf("failed to create FRP client instance: %w", err)
	}
	s.clientObj = svr

	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(context.Background())
	}

	go func() {
		s.status.SetStartedAt(s.timeSource.Get())
		s.status.SetIsRunning(true)
		defer func() {
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectly stopped, recovered from panic: %v\n%s", r, debug.Stack())
				s.handleInterrupt()
				_ = s.Stop()
			}
		}()

		if err := s.clientObj.Run(context.Background()); err != nil {
			logger.GetLogger(ID).Errorf("error running frpc: %v", err)
		}
		s.handleInterrupt()
	}()

	s.wg.Add(1)
	return nil
}

func (s *FrpClientServiceImpl) parseLocalServerAddr() (string, int, error) {
	hostname, port, err := net.SplitHostPort(s.localServerAddr)
	if err != nil {
		return "", 0, fmt.Errorf("failed to parse local server listening address: %w", err)
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return "", 0, fmt.Errorf("failed to parse local server listening port: %w", err)
	}

	return hostname, portNum, nil
}
