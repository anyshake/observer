package updater

import (
	"fmt"
)

func (s *UpdaterServiceImpl) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, con := range s.GetConfigConstraint() {
		if err := con.Init(s.actionHandler); err != nil {
			return fmt.Errorf("failed to initlize config constraint for service %s, namespace %s, key %s: %w", ID, con.GetNamespace(), con.GetKey(), err)
		}
	}

	releaseFetchUrl, err := (&updaterConfigReleaseFetchUrlImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.upgradeHelper.SetReleaseFetchUrl(releaseFetchUrl.(string))

	autoRestart, err := (&updaterConfigAutoRestartImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.autoRestartEnabled = autoRestart.(bool)

	return nil
}
