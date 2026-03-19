package mdns_discovery

import "fmt"

func (s *DiscoveryServiceImpl) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, con := range s.GetConfigConstraint() {
		if err := con.Init(s.actionHandler); err != nil {
			return fmt.Errorf("failed to initlize config constraint for service %s, namespace %s, key %s: %w", ID, con.GetNamespace(), con.GetKey(), err)
		}
	}

	instanceName, err := (&discoveryConfigInstanceNameImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.instanceName = instanceName.(string)

	return nil
}
