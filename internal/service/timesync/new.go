package timesync

import (
	"context"
	"time"

	"github.com/anyshake/observer/pkg/timesource"
)

func New(timeSource *timesource.Source) *TimeSyncServiceImpl {
	ctx, cancelFn := context.WithCancel(context.Background())
	obj := &TimeSyncServiceImpl{
		ctx:        ctx,
		cancelFn:   cancelFn,
		timeSource: timeSource,
	}
	obj.status.SetStartedAt(time.Unix(0, 0))
	obj.status.SetStoppedAt(time.Unix(0, 0))
	obj.status.SetUpdatedAt(time.Unix(0, 0))
	obj.status.SetIsRunning(false)
	obj.status.SetRestarts(0)
	return obj
}
