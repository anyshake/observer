package quakesense

func (s *QuakeSenseServiceImpl) IsEnabled() bool {
	enable, err := (&quakeSenseConfigEnabledImpl{}).Get(s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
