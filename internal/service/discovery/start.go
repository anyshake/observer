package mdns_discovery

import (
	"context"
	"fmt"
	"net"
	"runtime/debug"

	"github.com/anyshake/observer/pkg/logger"
	"github.com/grandcat/zeroconf"
)

func (s *DiscoveryServiceImpl) handleInterrupt() {
	s.wg.Done()
}

func (s *DiscoveryServiceImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(context.Background())
	}

	_, port, err := net.SplitHostPort(s.localServerAddr)
	if err != nil {
		return fmt.Errorf("invalid local server address: %w", err)
	}
	portInt, err := net.LookupPort("tcp", port)
	if err != nil {
		return fmt.Errorf("invalid local server port: %w", err)
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

		logger.GetLogger(ID).Infof("mDNS discovery service instance name: %s", s.instanceName)

		server, err := zeroconf.Register(s.instanceName, "_http._tcp", "local.", portInt, nil, nil)
		if err != nil {
			logger.GetLogger(ID).Errorf("failed to run mDNS discovery service: %v", err)
			s.handleInterrupt()
			_ = s.Stop()
			return
		}
		s.server = server

		<-s.ctx.Done()
		s.handleInterrupt()
	}()

	s.wg.Add(1)
	return nil
}
