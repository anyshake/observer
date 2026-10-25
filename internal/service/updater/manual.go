package updater

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/anyshake/observer/internal/notification"
	"github.com/anyshake/observer/pkg/semver"
)

var ErrUpgradeBusy = errors.New("an upgrade is already queued or running")

// RequestUpgrade queues work for the updater service. A stopped service is
// started for this request and its worker exits after the request completes.
func (s *UpdaterServiceImpl) RequestUpgrade(ctx context.Context) error {
	if !s.manualPending.CompareAndSwap(false, true) {
		return ErrUpgradeBusy
	}
	if err := ctx.Err(); err != nil {
		s.manualPending.Store(false)
		return err
	}
	if !s.GetStatus().GetIsRunning() {
		if err := s.Init(ctx); err != nil {
			s.manualPending.Store(false)
			return err
		}
		if err := s.start(ctx, true); err != nil {
			s.manualPending.Store(false)
			return err
		}
	}
	s.manualRequests <- struct{}{}
	return nil
}

func (s *UpdaterServiceImpl) notify(message string, level notification.Level) {
	if s.notifications != nil {
		s.notifications.Publish(notification.Event{
			ServiceID: ID,
			Message:   message,
			Level:     level,
			AdminOnly: true,
		})
	}
}

func (s *UpdaterServiceImpl) runManualUpgrade(ctx context.Context) error {
	latest, _, eligible, applied, err := s.upgradeHelper.CheckUpdateFresh(ctx)
	if err != nil {
		s.notify("Upgrade failed: version check failed: "+err.Error(), notification.LevelError)
		return err
	}
	if !eligible || applied {
		err := errors.New("no eligible, unapplied official upgrade is available")
		s.notify("Upgrade failed: "+err.Error(), notification.LevelError)
		return err
	}
	return s.applyUpdate(ctx, latest, true)
}

func (s *UpdaterServiceImpl) applyUpdate(ctx context.Context, latest *semver.Version, restart bool) error {
	s.notify("Upgrade to "+latest.String()+" started.", notification.LevelInfo)
	lastProgress := 0
	progress := func(percent int) {
		for next := lastProgress + 10; next <= percent; next += 10 {
			s.notify(fmt.Sprintf("Upgrade progress: %d%%", next), notification.LevelInfo)
			lastProgress = next
		}
	}
	fail := func(err error) error {
		s.notify("Upgrade failed: "+err.Error(), notification.LevelError)
		return err
	}

	progress(10)
	progress(20)
	release, _, err := s.upgradeHelper.FetchReleaseWithProgress(ctx, latest, UPDATE_FETCH_TIMEOUT, func(read, total int64) {
		progress(20 + int(50*read/total))
	})
	if err != nil {
		return fail(fmt.Errorf("release download or verification failed: %w", err))
	}
	progress(80)
	if err := ctx.Err(); err != nil {
		return fail(err)
	}
	progress(90)
	if err := s.upgradeHelper.ApplyUpgrade(latest, release); err != nil {
		return fail(fmt.Errorf("failed to apply release: %w", err))
	}
	progress(100)
	s.notify("Upgrade to "+latest.String()+" succeeded.", notification.LevelSuccess)
	if restart {
		s.notify("Application will restart in 3 seconds.", notification.LevelWarning)
		time.AfterFunc(3*time.Second, s.requestRestart)
	}
	return nil
}
