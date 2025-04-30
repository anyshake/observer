package forwarder

func (s *ForwarderServiceImpl) IsEnabled() bool {
	enable, err := (&forwarderConfigEnabledImpl{}).Get(s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
