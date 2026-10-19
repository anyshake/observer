package updater

import (
	"context"
	"fmt"
)

func (s *UpdaterServiceImpl) Init(ctx context.Context) error {
	if err := s.lifecycle.Lock(ctx); err != nil {
		return err
	}
	defer s.lifecycle.Unlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.GetIsRunning() {
		return fmt.Errorf("service %s cannot be initialized while running", ID)
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	for _, con := range s.GetConfigConstraint() {
		if err := con.Init(ctx, s.actionHandler); err != nil {
			return fmt.Errorf("failed to initlize config constraint for service %s, namespace %s, key %s: %w", ID, con.GetNamespace(), con.GetKey(), err)
		}
	}

	releaseFetchUrl, err := (&updaterConfigReleaseFetchUrlImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.upgradeHelper.SetReleaseFetchUrl(releaseFetchUrl.(string))

	autoRestart, err := (&updaterConfigAutoRestartImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.autoRestartEnabled = autoRestart.(bool)

	return nil
}
