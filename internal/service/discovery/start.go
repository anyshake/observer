package mdns_discovery

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime/debug"
	"time"

	service_forwarder "github.com/anyshake/observer/internal/service/forwarder"
	service_seedlink "github.com/anyshake/observer/internal/service/seedlink"
	service_winston "github.com/anyshake/observer/internal/service/winston"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/grandcat/zeroconf"
)

func (s *DiscoveryServiceImpl) newRegistration(instanceName, serviceType string, port int) (*mdnsRegistration, error) {
	srv, err := zeroconf.Register(instanceName, serviceType, "local.", port, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to register mDNS service %s on port %d: %w", serviceType, port, err)
	}
	return &mdnsRegistration{port: port, server: srv}, nil
}

func shutdownReg(reg *mdnsRegistration) {
	if reg != nil && reg.server != nil {
		reg.server.Shutdown()
	}
}

func (s *DiscoveryServiceImpl) syncRegistrations() {
	if s.httpReg == nil {
		_, port, err := net.SplitHostPort(s.webServerAddr)
		if err == nil {
			portInt, err := net.LookupPort("tcp", port)
			if err == nil {
				reg, err := s.newRegistration(s.instanceName, "_http._tcp", portInt)
				if err != nil {
					logger.GetLogger(ID).Errorf("failed to register HTTP mDNS service: %v", err)
				} else {
					s.httpReg = reg
					logger.GetLogger(ID).Infof("registered HTTP mDNS service on port %d", portInt)
				}
			}
		}
	}

	type candidate struct {
		enabled     bool
		serviceID   string
		serviceType string
		reg         **mdnsRegistration
		getPort     func() (int, bool) // port, isRunning
	}

	candidates := []candidate{
		{
			enabled:     s.registerSeedLink,
			serviceID:   service_seedlink.ID,
			serviceType: "_seedlink._tcp",
			reg:         &s.seedlinkReg,
			getPort: func() (int, bool) {
				svc, ok := s.serviceMap[service_seedlink.ID]
				if !ok {
					return 0, false
				}
				impl, ok := svc.(*service_seedlink.SeedLinkServiceImpl)
				if !ok {
					return 0, false
				}
				return impl.GetListenPort(), svc.GetStatus().GetIsRunning()
			},
		},
		{
			enabled:     s.registerWinston,
			serviceID:   service_winston.ID,
			serviceType: "_winston._tcp",
			reg:         &s.winstonReg,
			getPort: func() (int, bool) {
				svc, ok := s.serviceMap[service_winston.ID]
				if !ok {
					return 0, false
				}
				impl, ok := svc.(*service_winston.WinstonServiceImpl)
				if !ok {
					return 0, false
				}
				return impl.GetListenPort(), svc.GetStatus().GetIsRunning()
			},
		},
		{
			enabled:     s.registerForwarder,
			serviceID:   service_forwarder.ID,
			serviceType: "_forwarder._tcp",
			reg:         &s.forwarderReg,
			getPort: func() (int, bool) {
				svc, ok := s.serviceMap[service_forwarder.ID]
				if !ok {
					return 0, false
				}
				impl, ok := svc.(*service_forwarder.ForwarderServiceImpl)
				if !ok {
					return 0, false
				}
				return impl.GetListenPort(), svc.GetStatus().GetIsRunning()
			},
		},
	}

	for _, c := range candidates {
		reg := *c.reg

		if !c.enabled {
			if reg != nil {
				shutdownReg(reg)
				*c.reg = nil
				logger.GetLogger(ID).Infof("unregistered %s mDNS service (disabled)", c.serviceType)
			}
			continue
		}

		port, running := c.getPort()

		if !running {
			if reg != nil {
				shutdownReg(reg)
				*c.reg = nil
				logger.GetLogger(ID).Infof("unregistered %s mDNS service (service offline)", c.serviceType)
			}
			continue
		}

		if reg == nil || reg.port != port {
			if reg != nil {
				shutdownReg(reg)
				*c.reg = nil
			}
			newReg, err := s.newRegistration(s.instanceName, c.serviceType, port)
			if err != nil {
				logger.GetLogger(ID).Warnf("failed to register %s mDNS service: %v", c.serviceType, err)
			} else {
				*c.reg = newReg
				logger.GetLogger(ID).Infof("registered %s mDNS service on port %d", c.serviceType, port)
			}
		}
	}
}

func (s *DiscoveryServiceImpl) Start(ctx context.Context) error {
	if err := s.lifecycle.Lock(ctx); err != nil {
		return err
	}
	defer s.lifecycle.Unlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.GetIsRunning() {
		return errors.New("mDNS discovery service is already running")
	}
	if err := s.parentCtx.Err(); err != nil {
		return err
	}
	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(s.parentCtx)
	}

	runCtx, runCancel := s.ctx, s.cancelFn
	s.status.SetStartedAt(s.timeSource.Now())
	s.status.SetIsRunning(true)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(DISCOVERY_SCAN_INTERVAL)
		defer func() {
			ticker.Stop()
			runCancel()
			for _, reg := range []*mdnsRegistration{s.httpReg, s.seedlinkReg, s.winstonReg, s.forwarderReg} {
				shutdownReg(reg)
			}
			s.httpReg = nil
			s.seedlinkReg = nil
			s.winstonReg = nil
			s.forwarderReg = nil
			s.status.SetStoppedAt(s.timeSource.Now())
			s.status.SetIsRunning(false)
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectedly crashed, recovered from panic: %v\n%s", r, debug.Stack())
			}
		}()

		logger.GetLogger(ID).Infof("mDNS discovery service instance name: %s", s.instanceName)

		// Reset all registrations on (re)start
		s.httpReg = nil
		s.seedlinkReg = nil
		s.winstonReg = nil
		s.forwarderReg = nil

		// Run immediately, then on each tick
		s.syncRegistrations()

		for {
			select {
			case <-runCtx.Done():
				return
			case <-ticker.C:
				s.syncRegistrations()
			}
		}
	}()

	return nil
}
