package frp_client

import (
	"github.com/anyshake/observer/internal/service"
)

func (s *FrpClientServiceImpl) GetStatus() *service.Status {
	var status service.Status

	status.SetUpdatedAt(s.timeSource.Now())
	status.SetStoppedAt(s.status.GetStoppedAt())
	status.SetIsRunning(s.status.GetIsRunning())
	status.SetRestarts(s.status.GetRestarts())
	status.SetStartedAt(s.status.GetStartedAt())

	return &status
}
