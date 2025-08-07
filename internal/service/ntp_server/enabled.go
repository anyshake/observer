package ntp_server

func (s *NtpServerServiceImpl) IsEnabled() bool {
	enable, err := (&ntpServerConfigEnabledImpl{}).Get(s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
