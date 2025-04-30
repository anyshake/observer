package timesync

import (
	"context"
	"time"

	"github.com/anyshake/observer/pkg/logger"
)

func (s *TimeSyncServiceImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(context.Background())
	}

	go func() {
		s.status.SetStartedAt(s.timeSource.Get())
		s.status.SetIsRunning(true)

		timer := time.NewTimer(TIMESOURCE_REFRESH_INTERVAL)

		for {
			timer.Reset(TIMESOURCE_REFRESH_INTERVAL)

			select {
			case <-s.ctx.Done():
				timer.Stop()
				s.wg.Done()
				return
			case <-timer.C:
				if err := s.timeSource.Update(); err != nil {
					logger.GetLogger(ID).Warnf("failed to refresh time source: %v", err)
					continue
				}
				logger.GetLogger(ID).Debugf("time source refreshed: %v", s.timeSource.Get().Format(time.RFC3339))
			}
		}
	}()

	s.wg.Add(1)
	return nil
}
