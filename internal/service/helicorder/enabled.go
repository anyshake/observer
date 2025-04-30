package helicorder

func (s *HelicorderServiceImpl) IsEnabled() bool {
	enable, err := (&helicorderConfigEnabledImpl{}).Get(s.dataProvider.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); ok && !en {
		return false
	}
	return true
}
