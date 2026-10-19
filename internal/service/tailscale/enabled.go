package tailscale

import "context"

func (s *TailscaleServiceImpl) IsEnabled(ctx context.Context) bool {
	enabled, err := (&tailscaleConfigEnabledImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return false
	}
	if value, ok := enabled.(bool); !ok || !value {
		return false
	}
	return true
}
