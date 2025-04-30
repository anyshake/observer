package seedlink

func (s *SeedLinkServiceImpl) IsEnabled() bool {
	enable, err := (&seedlinkConfigEnabledImpl{}).Get(s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
