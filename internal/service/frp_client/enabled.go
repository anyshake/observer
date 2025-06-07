package frp_client

func (s *FrpClientServiceImpl) IsEnabled() bool {
	enable, err := (&frpClientConfigEnabledImpl{}).Get(s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
