package updater

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/anyshake/observer/pkg/logger"
)

func (s *UpdaterServiceImpl) handleInterrupt(ticker *time.Ticker) {
	ticker.Stop()
	s.wg.Done()
}

func (s *UpdaterServiceImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(context.Background())
	}

	go func() {
		ticker := time.NewTicker(UPDATE_CHECK_INTERVAL)

		s.status.SetStartedAt(s.timeSource.Now())
		s.status.SetIsRunning(true)
		defer func() {
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectly stopped, recovered from panic: %v\n%s", r, debug.Stack())
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
				latestVer, _, needUpdate, isApplied, err := s.upgradeHelper.CheckUpdate()
				if err != nil {
					logger.GetLogger(ID).Errorf("failed to check update: %v", err)
					continue
				}
				if needUpdate && !isApplied {
					logger.GetLogger(ID).Infof("found new version %s, upgrading process will start now", latestVer.String())
					release, url, err := s.upgradeHelper.FetchRelease(latestVer, UPDATE_FETCH_TIMEOUT)
					if err != nil {
						logger.GetLogger(ID).Errorf("failed to fetch release package: %v", err)
						continue
					}
					logger.GetLogger(ID).Infof("received %d bytes from URL: %s", len(release), url)
					if err := s.upgradeHelper.ApplyUpgrade(latestVer, release); err != nil {
						logger.GetLogger(ID).Errorf("failed to apply upgrade: %v", err)
						continue
					}
					logger.GetLogger(ID).Infof("new release %s has been applied successfully, restart to take effect", latestVer.String())
				} else if needUpdate && isApplied {
					logger.GetLogger(ID).Infoln("the latest version has already been applied, restart to take effect")
				} else {
					logger.GetLogger(ID).Infoln("no new version found for current software")
				}
			}
		}
	}()

	s.wg.Add(1)
	return nil
}
