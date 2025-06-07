package frp_client

import (
	"errors"
	"time"
)

func (s *FrpClientServiceImpl) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status.SetStoppedAt(s.timeSource.Get())
	s.status.SetIsRunning(false)
	s.cancelFn()

	done := make(chan struct{})
	go func() {
		s.clientObj.GracefulClose(500 * time.Millisecond)
		s.wg.Wait()
		close(done)
	}()

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case <-done:
		return nil
	case <-timer.C:
		return errors.New("timeout waiting for goroutines to finish")
	}
}
