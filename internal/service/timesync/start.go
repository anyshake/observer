package timesync

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/anyshake/observer/pkg/logger"
)

func (s *TimeSyncServiceImpl) handleInterrupt(ticker *time.Ticker) {
	ticker.Stop()
	s.wg.Done()
}

func (s *TimeSyncServiceImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(context.Background())
	}

	go func() {
		ticker := time.NewTicker(TIMESOURCE_REFRESH_INTERVAL)

		s.status.SetStartedAt(s.timeSource.Get())
		s.status.SetIsRunning(true)
		defer func() {
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectly crashed, recovered from panic: %v\n%s", r, debug.Stack())
				s.handleInterrupt(ticker)
				_ = s.Stop()
			}
		}()

		for {
			select {
			case <-s.ctx.Done():
				s.handleInterrupt(ticker)
				return
			case <-ticker.C:
				if err := s.timeSource.Update(); err != nil {
					logger.GetLogger(ID).Warnf("failed to refresh time source: %v", err)
				} else {
					logger.GetLogger(ID).Debugf("time source refreshed: %v", s.timeSource.Get().Format(time.RFC3339))
				}
			}
		}
	}()

	s.wg.Add(1)
	return nil
}
