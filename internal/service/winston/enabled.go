package winston

import "context"

func (s *WinstonServiceImpl) IsEnabled(ctx context.Context) bool {
	enable, err := (&winstonConfigEnabledImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
