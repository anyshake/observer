package updater

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/anyshake/observer/pkg/logger"
)

func (s *UpdaterServiceImpl) checkReadWritePermission() error {
	execDir := filepath.Dir(s.currentExePath)
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

func (s *UpdaterServiceImpl) Start(ctx context.Context) error {
	if err := s.lifecycle.Lock(ctx); err != nil {
		return err
	}
	defer s.lifecycle.Unlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.GetIsRunning() {
		return errors.New("updater service is already running")
	}
	if err := s.parentCtx.Err(); err != nil {
		return err
	}
	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(s.parentCtx)
	}

	if err := s.checkReadWritePermission(); err != nil {
		return fmt.Errorf("insufficient file system permissions for updater service: %w", err)
	}

	runCtx, runCancel := s.ctx, s.cancelFn
	s.status.SetStartedAt(s.timeSource.Now())
	s.status.SetIsRunning(true)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(UPDATE_CHECK_INTERVAL)
		defer func() {
			ticker.Stop()
			runCancel()
			s.status.SetStoppedAt(s.timeSource.Now())
			s.status.SetIsRunning(false)
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectedly crashed, recovered from panic: %v\n%s", r, debug.Stack())
			}
		}()

		for {
			select {
			case <-runCtx.Done():
				return
			case <-ticker.C:
				latestVer, _, needUpdate, isApplied, err := s.upgradeHelper.CheckUpdate(runCtx)
				if err != nil {
					if runCtx.Err() == nil {
						logger.GetLogger(ID).Errorf("failed to check update: %v", err)
					}
					continue
				}

				needRestart := false

				if needUpdate && !isApplied {
					logger.GetLogger(ID).Infof("found new version %s, upgrading process will start now", latestVer.String())
					release, url, err := s.upgradeHelper.FetchRelease(runCtx, latestVer, UPDATE_FETCH_TIMEOUT)
					if err != nil {
						if runCtx.Err() == nil {
							logger.GetLogger(ID).Errorf("failed to fetch release package: %v", err)
						}
						continue
					}
					logger.GetLogger(ID).Infof("received %d bytes from URL: %s", len(release), url)
					if err := s.upgradeHelper.ApplyUpgrade(latestVer, release); err != nil {
						logger.GetLogger(ID).Errorf("failed to apply upgrade: %v", err)
						continue
					}
					if !s.autoRestartEnabled {
						logger.GetLogger(ID).Infof("new release %s has been applied successfully, restart to take effect", latestVer.String())
					}
					needRestart = s.autoRestartEnabled
				} else if needUpdate && isApplied {
					if !s.autoRestartEnabled {
						logger.GetLogger(ID).Infoln("the new version has already been applied, restart to take effect")
					}
					needRestart = s.autoRestartEnabled
				} else {
					logger.GetLogger(ID).Infoln("no new version found for current software")
					needRestart = false
				}

				if needRestart {
					logger.GetLogger(ID).Infoln("restarting application in 3 second...")
					time.AfterFunc(3*time.Second, s.requestRestart)
				}
			}
		}
	}()

	return nil
}

func (s *UpdaterServiceImpl) requestRestart() {
	select {
	case s.restartChan <- struct{}{}:
	default:
	}
}
