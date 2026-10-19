package quakesense

import (
	"context"
)

func (s *QuakeSenseServiceImpl) Stop(ctx context.Context) error {
	return s.lifecycle.Stop(ctx, func() error {
		s.mu.Lock()
		if s.cancelFn != nil {
			s.cancelFn()
		}
		s.status.SetIsRunning(false)
		s.mu.Unlock()
		defer func() {
			s.status.SetStoppedAt(s.timeSource.Now())
			s.status.SetIsRunning(false)
		}()

		_ = s.hardwareDev.UnsubscribeRealtime(ID)
		s.wg.Wait()
		return nil
	})
}
