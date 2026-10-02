package tailscale

import "fmt"

func (s *TailscaleServiceImpl) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, con := range s.GetConfigConstraint() {
		if err := con.Init(s.actionHandler); err != nil {
			return fmt.Errorf("failed to initialize config constraint for service %s, namespace %s, key %s: %w", ID, con.GetNamespace(), con.GetKey(), err)
		}
	}

	hostname, err := (&tailscaleConfigHostnameImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	controlURL, err := (&tailscaleConfigControlURLImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	forwardPorts, err := (&tailscaleConfigForwardPortsImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	authKey, err := (&tailscaleConfigAuthKeyImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	stateStore, err := s.newStateStore(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to initialize Tailscale state store: %w", err)
	}

	s.hostname = hostname.(string)
	s.controlURL = controlURL.(string)
	s.forwardPorts = forwardPorts.([]string)
	s.authKey = authKey.(string)
	s.stateStore = stateStore
	return nil
}
