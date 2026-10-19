package metrics

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/host"
)

func (s *MetricsServiceImpl) Init(ctx context.Context) error {
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

	hostInfo, err := host.InfoWithContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get host info: %w", err)
	}
	s.hostInfo = hostInfo

	if s.hostInfo.PlatformVersion == "" {
		s.hostInfo.PlatformVersion = ATTRIBUTE_DUMMY_VALUE
	}

	if s.hostInfo.Platform == "" {
		s.hostInfo.Platform = ATTRIBUTE_DUMMY_VALUE
	}

	if s.hostInfo.Hostname == "" {
		s.hostInfo.Hostname = ATTRIBUTE_DUMMY_VALUE
	}

	s.startTime = s.timeSource.Now()

	return nil
}
