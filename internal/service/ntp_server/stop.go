package ntp_server

import (
	"context"
)

func (s *NtpServerServiceImpl) Stop(ctx context.Context) error {
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
		return nil
	})
}
