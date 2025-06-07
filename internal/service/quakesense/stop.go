package quakesense

import (
	"errors"
	"time"
)

func (s *QuakeSenseServiceImpl) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status.SetStoppedAt(s.timeSource.Get())
	s.status.SetIsRunning(false)
	s.cancelFn()

	done := make(chan struct{})
	go func() {
		s.prevSamplerate = 0
		s.filterKernel = nil
		_ = s.hardwareDev.Unsubscribe(ID)
		if s.mqttClient != nil {
			s.mqttClient.Disconnect(100)
		}
		s.channelBuffer.Reset()
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
