package watchcat

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/anyshake/observer/pkg/logger"
)

func (s *WatchCatServiceImpl) handleInterrupt(ticker *time.Ticker) {
	ticker.Stop()
	s.wg.Done()
}

func (s *WatchCatServiceImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(context.Background())
	}

	go func() {
		ticker := time.NewTicker(WATCHCAT_CHECK_INTERVAL)

		s.status.SetStartedAt(s.timeSource.Get())
		s.status.SetIsRunning(true)
		defer func() {
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectly stopped, recovered from panic: %v\n%s", r, debug.Stack())
				s.handleInterrupt(ticker)
				_ = s.Stop()
			}
		}()

		var lastUpdatedAt time.Time
		for {
			select {
			case <-s.ctx.Done():
				s.handleInterrupt(ticker)
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

	s.wg.Add(1)
	return nil
}
