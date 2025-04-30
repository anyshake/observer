package seedlink

import "github.com/anyshake/observer/internal/service"

func (s *SeedLinkServiceImpl) GetStatus() *service.Status {
	var status service.Status

	status.SetUpdatedAt(s.timeSource.Get())
	status.SetStoppedAt(s.status.GetStoppedAt())
	status.SetIsRunning(s.status.GetIsRunning())
	status.SetRestarts(s.status.GetRestarts())
	status.SetStartedAt(s.status.GetStartedAt())

	return &status
}
