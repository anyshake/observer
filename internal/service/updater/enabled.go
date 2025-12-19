package updater

func (s *UpdaterServiceImpl) IsEnabled() bool {
	enable, err := (&updaterConfigEnabledImpl{}).Get(s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
