package frp_client

import (
	"errors"
	"time"
)

func (s *FrpClientServiceImpl) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status.SetStoppedAt(s.timeSource.Now())
	s.status.SetIsRunning(false)
	s.cancelFn()

	done := make(chan struct{})
	go func() {
		// Run observes s.ctx and returns before Close reads the client.
		s.wg.Wait()
		if s.clientObj != nil {
			s.clientObj.Close()
		}
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
