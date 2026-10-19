package tailscale

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"

	"tailscale.com/ipn"
)

func hasPersistedIdentity(store ipn.StateStore) (bool, error) {
	if store == nil {
		return false, errors.New("tailscale state store is not configured")
	}

	identityKeys := []ipn.StateKey{
		ipn.CurrentProfileStateKey,
		ipn.ServerModeStartKey,
		ipn.LegacyGlobalDaemonStateKey,
	}
	for _, key := range identityKeys {
		value, err := store.ReadState(key)
		switch {
		case err == nil && len(value) > 0:
			return true, nil
		case err == nil:
			continue
		case errors.Is(err, ipn.ErrStateNotExist):
			continue
		default:
			return false, fmt.Errorf("failed to inspect persisted tailscale identity: %w", err)
		}
	}

	return false, nil
}

func (s *TailscaleServiceImpl) redact(value string) string {
	if s.authKey == "" {
		return value
	}
	return strings.ReplaceAll(value, s.authKey, "[redacted]")
}

func (s *TailscaleServiceImpl) userLogf(format string, args ...any) {
	message := s.redact(fmt.Sprintf(format, args...))
	s.log.Infof("tsnet: %s", message)
}

func (s *TailscaleServiceImpl) Start(ctx context.Context) error {
	if err := s.lifecycle.Lock(ctx); err != nil {
		return err
	}
	defer s.lifecycle.Unlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.GetIsRunning() {
		return errors.New("tailscale service is already running")
	}
	if err := s.parentCtx.Err(); err != nil {
		return err
	}
	forwards, err := buildForwards(s.webServerAddr, s.forwardPorts)
	if err != nil {
		return err
	}

	hasIdentity, err := hasPersistedIdentity(s.stateStore)
	if err != nil {
		return err
	}
	if !hasIdentity && strings.TrimSpace(s.authKey) == "" {
		return errors.New("tailscale auth key is required for initial registration")
	}

	s.ctx, s.cancelFn = context.WithCancel(s.parentCtx)
	runCtx, runCancel := s.ctx, s.cancelFn
	authKey := s.authKey
	if hasIdentity {
		authKey = ""
	}

	node := s.newNode(nodeConfig{
		Hostname:   s.hostname,
		AuthKey:    authKey,
		ControlURL: s.controlURL,
		Store:      s.stateStore,
		UserLogf:   s.userLogf,
	})
	if err := node.Start(); err != nil {
		runCancel()
		return fmt.Errorf("failed to start embedded tailscale node: %s", s.redact(err.Error()))
	}
	s.status.SetStartedAt(s.timeSource.Now())
	s.status.SetIsRunning(true)

	s.wg.Add(1)
	go s.run(runCtx, runCancel, node, forwards)
	return nil
}

func (s *TailscaleServiceImpl) run(ctx context.Context, cancel context.CancelFunc, node embeddedNode, forwards proxyForwards) {
	defer s.wg.Done()
	defer func() {
		cancel()
		if err := node.Close(); err != nil {
			s.log.Errorf("failed to close embedded tailscale node: %s", s.redact(err.Error()))
		}
		s.status.SetStoppedAt(s.timeSource.Now())
		s.status.SetIsRunning(false)
		if r := recover(); r != nil {
			s.log.Errorf("service unexpectedly crashed, recovered from panic: %v\n%s", r, debug.Stack())
		}
	}()

	upErr := node.Up(ctx)
	if upErr == nil {
		if err := serveForwards(ctx, node, forwards, s.log); err != nil && ctx.Err() == nil {
			s.log.Errorf("embedded tailscale proxy failed: %s", s.redact(err.Error()))
		}
	} else if ctx.Err() == nil {
		s.log.Errorf("embedded tailscale node stopped: %s", s.redact(upErr.Error()))
	}
}
