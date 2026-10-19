package miniseed

import "context"

func (s *MiniSeedServiceImpl) IsEnabled(ctx context.Context) bool {
	enable, err := (&miniSeedConfigEnabledImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
