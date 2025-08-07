package ntp_server

import "fmt"

func (s *NtpServerServiceImpl) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, con := range s.GetConfigConstraint() {
		if err := con.Init(s.actionHandler); err != nil {
			return fmt.Errorf("failed to initlize config constraint for service %s, namespace %s, key %s: %w", ID, con.GetNamespace(), con.GetKey(), err)
		}
	}

	listenHost, err := (&ntpServerConfigListenHostImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.listenHost = listenHost.(string)

	listenPort, err := (&ntpServerConfigListenPortImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.listenPort = listenPort.(int)

	return nil
}
