package setup_station

import (
	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
)

type SetupStationStartupImpl struct {
	ActionHandler            *action.Handler
	StationConfigConstraints []config.IConstraint
}
