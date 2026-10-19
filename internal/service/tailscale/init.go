package tailscale

import (
	"context"
	"fmt"
)

func (s *TailscaleServiceImpl) Init(ctx context.Context) error {
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
			return fmt.Errorf("failed to initialize config constraint for service %s, namespace %s, key %s: %w", ID, con.GetNamespace(), con.GetKey(), err)
		}
	}

	hostname, err := (&tailscaleConfigHostnameImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	controlURL, err := (&tailscaleConfigControlURLImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	forwardPorts, err := (&tailscaleConfigForwardPortsImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	authKey, err := (&tailscaleConfigAuthKeyImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	stateStore, err := s.newStateStore(ctx, s.actionHandler)
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
