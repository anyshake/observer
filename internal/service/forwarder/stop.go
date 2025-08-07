package forwarder

import (
	"errors"
	"time"
)

func (s *ForwarderServiceImpl) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status.SetStoppedAt(s.timeSource.Now())
	s.status.SetIsRunning(false)
	s.cancelFn()

	done := make(chan struct{})
	go func() {
		_ = s.hardwareDev.Unsubscribe(ID)
		if s.listener != nil {
			_ = s.listener.Close()
		}
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
