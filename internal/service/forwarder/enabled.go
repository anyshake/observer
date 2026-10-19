package forwarder

import "context"

func (s *ForwarderServiceImpl) IsEnabled(ctx context.Context) bool {
	enable, err := (&forwarderConfigEnabledImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
