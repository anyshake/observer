package tailscale

import (
	"context"
	"errors"
	"fmt"
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

func (s *TailscaleServiceImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.GetIsRunning() {
		return errors.New("tailscale service is already running")
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

	s.ctx, s.cancelFn = context.WithCancel(context.Background())
	s.setRunError(nil)
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
		s.cancelFn()
		return fmt.Errorf("failed to start embedded tailscale node: %s", s.redact(err.Error()))
	}
	s.status.SetStartedAt(s.timeSource.Now())
	s.status.SetIsRunning(true)

	s.wg.Add(1)
	go s.run(node, s.ctx, forwards)
	return nil
}

func (s *TailscaleServiceImpl) run(node embeddedNode, ctx context.Context, forwards proxyForwards) {
	defer s.wg.Done()

	upErr := node.Up(ctx)
	if upErr == nil {
		if err := serveForwards(ctx, node, forwards, s.log); err != nil && ctx.Err() == nil {
			s.setRunError(fmt.Errorf("embedded tailscale proxy failed: %s", s.redact(err.Error())))
			s.log.Errorf("embedded tailscale proxy failed: %s", s.redact(err.Error()))
		}
	} else if ctx.Err() == nil {
		s.log.Errorf("embedded tailscale node stopped: %s", s.redact(upErr.Error()))
	}

	closeErr := node.Close()
	if closeErr != nil {
		closeErr = fmt.Errorf("failed to close embedded tailscale node: %s", s.redact(closeErr.Error()))
		s.log.Errorf("%v", closeErr)
	}
	if closeErr != nil || s.getRunError() == nil {
		s.setRunError(closeErr)
	}
	s.status.SetStoppedAt(s.timeSource.Now())
	s.status.SetIsRunning(false)
}

func (s *TailscaleServiceImpl) setRunError(err error) {
	s.errMu.Lock()
	s.runError = err
	s.errMu.Unlock()
}

func (s *TailscaleServiceImpl) getRunError() error {
	s.errMu.Lock()
	defer s.errMu.Unlock()
	return s.runError
}
