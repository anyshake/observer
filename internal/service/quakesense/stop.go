package quakesense

import (
	"errors"
	"time"
)

func (s *QuakeSenseServiceImpl) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status.SetStoppedAt(s.timeSource.Now())
	s.status.SetIsRunning(false)
	s.cancelFn()

	done := make(chan struct{})
	go func() {
		_ = s.hardwareDev.UnsubscribeRealtime(ID)
		s.wg.Wait()
		s.prevSamplerate = 0
		s.filterKernel = nil
		s.mqttClient = nil
		if s.channelBuffer != nil {
			s.channelBuffer.Reset()
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
