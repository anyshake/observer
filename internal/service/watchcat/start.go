package watchcat

import (
	"context"
	"errors"
	"runtime/debug"
	"time"

	"github.com/anyshake/observer/pkg/logger"
)

func (s *WatchCatServiceImpl) Start(ctx context.Context) error {
	if err := s.lifecycle.Lock(ctx); err != nil {
		return err
	}
	defer s.lifecycle.Unlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.GetIsRunning() {
		return errors.New("WatchCat service is already running")
	}
	if err := s.parentCtx.Err(); err != nil {
		return err
	}
	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(s.parentCtx)
	}

	runCtx, runCancel := s.ctx, s.cancelFn
	s.status.SetStartedAt(s.timeSource.Now())
	s.status.SetIsRunning(true)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(WATCHCAT_CHECK_INTERVAL)
		defer func() {
			ticker.Stop()
			runCancel()
			s.status.SetStoppedAt(s.timeSource.Now())
			s.status.SetIsRunning(false)
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectedly crashed, recovered from panic: %v\n%s", r, debug.Stack())
			}
		}()

		var lastUpdatedAt time.Time
		for {
			select {
			case <-runCtx.Done():
				return
			case <-ticker.C:
				status := s.hardwareDev.GetStatus()
				updatedAt := status.GetUpdatedAt()
				if !lastUpdatedAt.IsZero() && updatedAt.Sub(lastUpdatedAt) == 0 {
					logger.GetLogger(ID).Warnf("device is not responsding as expected, trying to flush device")
					if err := s.hardwareDev.Flush(); err != nil {
						logger.GetLogger(ID).Errorf("failed to flush device: %v", err)
					}
				}
				lastUpdatedAt = updatedAt
			}
		}
	}()

	return nil
}
