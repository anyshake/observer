package updater

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/anyshake/observer/pkg/logger"
)

func (s *UpdaterServiceImpl) checkReadWritePermission() error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	execDir := filepath.Dir(execPath)
	testFile := filepath.Join(execDir, ".observer_updater_write_test")
	testContent := fmt.Appendf(nil, "anyshake-observer-updater-service-%d", time.Now().UnixNano())

	if err := os.WriteFile(testFile, testContent, 0o644); err != nil {
		return fmt.Errorf("no write permission on executable dir %s: %w", execDir, err)
	}

	readBack, err := os.ReadFile(testFile)
	if err != nil {
		_ = os.Remove(testFile)
		return fmt.Errorf("failed to read test file in executable dir %s: %w", execDir, err)
	}

	if !bytes.Equal(readBack, testContent) {
		_ = os.Remove(testFile)
		return fmt.Errorf("content mismatch when testing write permission in %s", execDir)
	}

	if err := os.Remove(testFile); err != nil {
		return fmt.Errorf("failed to cleanup test file in %s: %w", execDir, err)
	}

	return nil
}

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

	if err := s.checkReadWritePermission(); err != nil {
		return fmt.Errorf("insufficient file system permissions for updater service: %w", err)
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
