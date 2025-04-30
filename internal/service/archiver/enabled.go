package archiver

func (s *ArchiverServiceImpl) IsEnabled() bool {
	enable, err := (&archiverConfigEnabledImpl{}).Get(s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
