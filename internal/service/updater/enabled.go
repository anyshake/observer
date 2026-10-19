package updater

import "context"

func (s *UpdaterServiceImpl) IsEnabled(ctx context.Context) bool {
	enable, err := (&updaterConfigEnabledImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
