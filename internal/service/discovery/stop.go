package mdns_discovery

import (
	"errors"
	"time"
)

func (s *DiscoveryServiceImpl) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status.SetStoppedAt(s.timeSource.Now())
	s.status.SetIsRunning(false)
	s.cancelFn()

	done := make(chan struct{})
	go func() {
		// Wait until the scan loop has finished writing the registration
		// pointers, then shut those servers down.
		s.wg.Wait()
		for _, reg := range []*mdnsRegistration{
			s.httpReg, s.seedlinkReg, s.winstonReg, s.forwarderReg,
		} {
			if reg != nil && reg.server != nil {
				reg.server.Shutdown()
			}
		}
		close(done)
	}()

	timer := time.NewTimer(3 * time.Second)
	defer timer.Stop()

	select {
	case <-done:
		return nil
	case <-timer.C:
		return errors.New("timeout waiting for goroutines to finish")
	}
}
