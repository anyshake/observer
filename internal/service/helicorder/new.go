package helicorder

import (
	"context"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/timesource"
)

func New(hardwareDev hardware.IHardware, actionHandler *action.Handler, timeSource *timesource.Source) *HelicorderServiceImpl {
	ctx, cancelFn := context.WithCancel(context.Background())
	obj := &HelicorderServiceImpl{
		ctx:         ctx,
		cancelFn:    cancelFn,
		timeSource:  timeSource,
		hardwareDev: hardwareDev,
		dataProvider: provider{
			actionHandler: actionHandler,
			queryCache:    cache.NewKv[[]queryCacheData](time.Hour),
		},
	}
	obj.status.SetStartedAt(time.Unix(0, 0))
	obj.status.SetStoppedAt(time.Unix(0, 0))
	obj.status.SetUpdatedAt(time.Unix(0, 0))
	obj.status.SetIsRunning(false)
	obj.status.SetRestarts(0)
	return obj
}
