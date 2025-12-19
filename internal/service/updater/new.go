package updater

import (
	"context"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/upgrade"
	"github.com/anyshake/observer/pkg/timesource"
)

func New(actionHandler *action.Handler, timeSource *timesource.Source, helper *upgrade.Helper) *UpdaterServiceImpl {
	ctx, cancelFn := context.WithCancel(context.Background())
	obj := &UpdaterServiceImpl{
		ctx:           ctx,
		cancelFn:      cancelFn,
		upgradeHelper: helper,
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
