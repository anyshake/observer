package helicorder

import (
	"context"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/pkg/timesource"
)

func New(hardwareDev hardware.IHardware, actionHandler *action.Handler, timeSource *timesource.Source) *HelicorderServiceImpl {
	ctx, cancelFn := context.WithCancel(context.Background())
	obj := &HelicorderServiceImpl{
		ctx:          ctx,
		cancelFn:     cancelFn,
		timeSource:   timeSource,
		hardwareDev:  hardwareDev,
		cacheStorage: CACHE_STORAGE_DISABLED,
		cachePath:    CACHE_DEFAULT_PATH,
		dataProvider: provider{
			ctx:           ctx,
			actionHandler: actionHandler,
		},
	}
	obj.status.SetStartedAt(time.Unix(0, 0))
	obj.status.SetStoppedAt(time.Unix(0, 0))
	obj.status.SetUpdatedAt(time.Unix(0, 0))
	obj.status.SetIsRunning(false)
	obj.status.SetRestarts(0)
	return obj
}
