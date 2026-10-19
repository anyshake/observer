package ntp_server

import (
	"context"
	"fmt"
)

func (s *NtpServerServiceImpl) Init(ctx context.Context) error {
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

	listenHost, err := (&ntpServerConfigListenHostImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.listenHost = listenHost.(string)

	listenPort, err := (&ntpServerConfigListenPortImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.listenPort = listenPort.(int)

	delayControl, err := (&ntpServerConfigDelayControlImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.delayControl = delayControl.(int64)

	return nil
}
