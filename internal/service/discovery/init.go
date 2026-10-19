package mdns_discovery

import (
	"context"
	"fmt"
)

func (s *DiscoveryServiceImpl) Init(ctx context.Context) error {
	if err := s.lifecycle.Lock(ctx); err != nil {
		return err
	}
	defer s.lifecycle.Unlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.GetIsRunning() {
		return fmt.Errorf("service %s cannot be initialized while running", ID)
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	for _, con := range s.GetConfigConstraint() {
		if err := con.Init(ctx, s.actionHandler); err != nil {
			return fmt.Errorf("failed to initlize config constraint for service %s, namespace %s, key %s: %w", ID, con.GetNamespace(), con.GetKey(), err)
		}
	}

	instanceName, err := (&discoveryConfigInstanceNameImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.instanceName = instanceName.(string)

	registerSeedLink, err := (&discoveryConfigRegisterSeedLinkImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.registerSeedLink = registerSeedLink.(bool)

	registerWinston, err := (&discoveryConfigRegisterWinstonImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.registerWinston = registerWinston.(bool)

	registerForwarder, err := (&discoveryConfigRegisterForwarderImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.registerForwarder = registerForwarder.(bool)

	return nil
}
