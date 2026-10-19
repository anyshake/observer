package helicorder

import "context"

func (s *HelicorderServiceImpl) IsEnabled(ctx context.Context) bool {
	enable, err := (&helicorderConfigEnabledImpl{}).Get(ctx, s.dataProvider.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); ok && !en {
		return false
	}
	return true
}
