package updater

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/anyshake/observer/internal/notification"
	"github.com/anyshake/observer/pkg/logger"
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
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ctx.Err(); err != nil {
		s.manualPending.Store(false)
		return err
	}
	requestCtx, cancel := context.WithCancel(s.ctx)
	s.manualCancel = cancel
	s.manualApplying = false
	s.notify("Manual upgrade queued.", notification.LevelInfo)
	s.manualRequests <- requestCtx
	return nil
}

// CancelUpgrade interrupts checking and downloading. Executable replacement is
// atomic from the user's perspective and must finish once it has started.
func (s *UpdaterServiceImpl) CancelUpgrade() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.manualCancel == nil || s.manualApplying {
		return errors.New("no cancellable manual upgrade is running")
	}
	s.manualCancel()
	return nil
}

func (s *UpdaterServiceImpl) ManualUpgradeStatus() (running, cancellable bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.manualPending.Load(), s.manualCancel != nil && !s.manualApplying
}

func (s *UpdaterServiceImpl) finishManualUpgrade(manualOnly bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.manualCancel != nil {
		s.manualCancel()
		s.manualCancel = nil
	}
	s.manualApplying = false
	if !manualOnly {
		s.manualPending.Store(false)
	}
}

func (s *UpdaterServiceImpl) notify(message string, level notification.Level) {
	log := logger.GetLogger(ID)
	logMessage := lowerFirst(message)
	switch level {
	case notification.LevelError:
		log.Errorln(logMessage)
	case notification.LevelWarning:
		log.Warnln(logMessage)
	default:
		log.Infoln(logMessage)
	}
	if s.notifications != nil {
		s.notifications.Publish(notification.Event{
			ServiceID: ID,
			Message:   message,
			Level:     level,
			AdminOnly: true,
		})
	}
}

func (s *UpdaterServiceImpl) upgradeFailed(ctx context.Context, err error) error {
	if errors.Is(ctx.Err(), context.Canceled) {
		s.notify("Upgrade cancelled.", notification.LevelWarning)
		return context.Canceled
	}
	s.notify("Upgrade failed: "+err.Error(), notification.LevelError)
	return err
}

func (s *UpdaterServiceImpl) runManualUpgrade(ctx context.Context) error {
	s.notify("Checking for an official upgrade.", notification.LevelInfo)
	if err := ctx.Err(); err != nil {
		return s.upgradeFailed(ctx, err)
	}
	latest, _, eligible, applied, err := s.upgradeHelper.CheckUpdateFresh(ctx)
	if err != nil {
		return s.upgradeFailed(ctx, fmt.Errorf("version check failed: %w", err))
	}
	if !eligible || applied {
		err := errors.New("no eligible, unapplied official upgrade is available")
		return s.upgradeFailed(ctx, err)
	}
	return s.applyUpdate(ctx, latest, true, true)
}

func (s *UpdaterServiceImpl) applyUpdate(ctx context.Context, latest *semver.Version, restart, manual bool) error {
	if err := ctx.Err(); err != nil {
		return s.upgradeFailed(ctx, err)
	}
	s.notify("Upgrade to "+latest.String()+" started.", notification.LevelInfo)
	lastProgress := 0
	progress := func(percent int) {
		for next := lastProgress + 10; next <= percent; next += 10 {
			s.notify(fmt.Sprintf("Upgrade progress: %d%%", next), notification.LevelInfo)
			lastProgress = next
		}
	}
	fail := func(err error) error {
		return s.upgradeFailed(ctx, err)
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
	s.mu.Lock()
	err = ctx.Err()
	if err == nil && manual {
		s.manualApplying = true
	}
	s.mu.Unlock()
	if err != nil {
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

// lowerFirst keeps log lines in the lower-case style used elsewhere while
// notifications shown to users remain sentence-cased.
func lowerFirst(message string) string {
	r, size := utf8.DecodeRuneInString(message)
	if r == utf8.RuneError || !unicode.IsUpper(r) {
		return message
	}
	return string(unicode.ToLower(r)) + message[size:]
}
