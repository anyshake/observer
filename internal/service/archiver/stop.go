package archiver

import (
	"context"
)

func (s *ArchiverServiceImpl) Stop(ctx context.Context) error {
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

		_ = s.hardwareDev.Unsubscribe(ID)
		s.wg.Wait()
		return nil
	})
}
