package tailscale

func (s *TailscaleServiceImpl) IsEnabled() bool {
	enabled, err := (&tailscaleConfigEnabledImpl{}).Get(s.actionHandler)
	if err != nil {
		return false
	}
	if value, ok := enabled.(bool); !ok || !value {
		return false
	}
	return true
}
