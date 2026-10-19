package quakesense

import "context"

func (s *QuakeSenseServiceImpl) IsEnabled(ctx context.Context) bool {
	enable, err := (&quakeSenseConfigEnabledImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
