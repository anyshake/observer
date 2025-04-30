package miniseed

func (s *MiniSeedServiceImpl) IsEnabled() bool {
	enable, err := (&miniSeedConfigEnabledImpl{}).Get(s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
