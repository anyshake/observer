package ntp_server

import "context"

func (s *NtpServerServiceImpl) IsEnabled(ctx context.Context) bool {
	enable, err := (&ntpServerConfigEnabledImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
