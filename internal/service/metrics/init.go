package metrics

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/host"
)

func (s *MetricsServiceImpl) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, con := range s.GetConfigConstraint() {
		if err := con.Init(s.actionHandler); err != nil {
			return fmt.Errorf("failed to initlize config constraint for service %s, namespace %s, key %s: %w", ID, con.GetNamespace(), con.GetKey(), err)
		}
	}

	hostInfo, err := host.Info()
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

	s.startTime = s.timeSource.Get()

	return nil
}
