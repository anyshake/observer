package metrics

import (
	"context"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/pkg/semver"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/anyshake/observer/pkg/unibuild"
)

func New(hardwareDev hardware.IHardware, actionHandler *action.Handler, timeSource *timesource.Source, version *semver.Version, build *unibuild.UniBuild) *MetricsServiceImpl {
	ctx, cancelFn := context.WithCancel(context.Background())
	obj := &MetricsServiceImpl{
		ctx:      ctx,
		cancelFn: cancelFn,

		actionHandler: actionHandler,
		timeSource:    timeSource,
		hardwareDev:   hardwareDev,

		version: version,
		build:   build,
	}
	obj.status.SetStartedAt(time.Unix(0, 0))
	obj.status.SetStoppedAt(time.Unix(0, 0))
	obj.status.SetUpdatedAt(time.Unix(0, 0))
	obj.status.SetIsRunning(false)
	obj.status.SetRestarts(0)
	return obj
}
