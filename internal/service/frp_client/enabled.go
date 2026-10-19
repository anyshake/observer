package frp_client

import "context"

func (s *FrpClientServiceImpl) IsEnabled(ctx context.Context) bool {
	enable, err := (&frpClientConfigEnabledImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
