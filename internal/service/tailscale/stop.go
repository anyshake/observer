package tailscale

import (
	"errors"
	"time"
)

func (s *TailscaleServiceImpl) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancelFn != nil {
		s.cancelFn()
	}
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	timer := time.NewTimer(shutdownTimeout)
	defer timer.Stop()

	select {
	case <-done:
		return s.getRunError()
	case <-timer.C:
		return errors.New("timeout waiting for embedded tailscale node to stop")
	}
}
