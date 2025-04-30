package metrics

import (
	"context"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/pkg/timesource"
)

func New(hardwareDev hardware.IHardware, actionHandler *action.Handler, timeSource *timesource.Source, binaryVersion, commitHash, buildPlatform string) *MetricsServiceImpl {
	ctx, cancelFn := context.WithCancel(context.Background())
	obj := &MetricsServiceImpl{
		ctx:      ctx,
		cancelFn: cancelFn,

		actionHandler: actionHandler,
		timeSource:    timeSource,
		hardwareDev:   hardwareDev,

		binaryVersion: binaryVersion,
		commitHash:    commitHash,
		buildPlatform: buildPlatform,
	}
	obj.status.SetStartedAt(time.Unix(0, 0))
	obj.status.SetStoppedAt(time.Unix(0, 0))
	obj.status.SetUpdatedAt(time.Unix(0, 0))
	obj.status.SetIsRunning(false)
	obj.status.SetRestarts(0)
	return obj
}
