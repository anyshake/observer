package helicorder

import (
	"context"
	"fmt"
)

func (s *HelicorderServiceImpl) Stop(ctx context.Context) error {
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

		s.wg.Wait()
		if s.dataProvider.queryCache != nil {
			if err := s.dataProvider.queryCache.Clear(); err != nil {
				return fmt.Errorf("failed to clear helicorder cache: %w", err)
			}
		}
		return nil
	})
}
