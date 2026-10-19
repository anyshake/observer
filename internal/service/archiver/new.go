package archiver

import (
	"context"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/pkg/timesource"
)

func New(parentCtx context.Context, hardwareDev hardware.IHardware, actionHandler *action.Handler, timeSource *timesource.Source) *ArchiverServiceImpl {
	ctx, cancelFn := context.WithCancel(parentCtx)
	obj := &ArchiverServiceImpl{
		hardwareDev:   hardwareDev,
		parentCtx:     parentCtx,
		ctx:           ctx,
		cancelFn:      cancelFn,
		timeSource:    timeSource,
		actionHandler: actionHandler,
	}
	obj.status.SetStartedAt(time.Unix(0, 0))
	obj.status.SetStoppedAt(time.Unix(0, 0))
	obj.status.SetUpdatedAt(time.Unix(0, 0))
	obj.status.SetIsRunning(false)
	obj.status.SetRestarts(0)
	return obj
}
