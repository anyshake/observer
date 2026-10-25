package updater

import (
	"context"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/notification"
	"github.com/anyshake/observer/internal/upgrade"
	"github.com/anyshake/observer/pkg/timesource"
)

func New(parentCtx context.Context, actionHandler *action.Handler, timeSource *timesource.Source, helper *upgrade.Helper, currentExePath string, restartChan chan struct{}, notifications notification.Publisher) *UpdaterServiceImpl {
	ctx, cancelFn := context.WithCancel(parentCtx)
	obj := &UpdaterServiceImpl{
		parentCtx:      parentCtx,
		ctx:            ctx,
		cancelFn:       cancelFn,
		upgradeHelper:  helper,
		timeSource:     timeSource,
		actionHandler:  actionHandler,
		currentExePath: currentExePath,
		restartChan:    restartChan,
		notifications:  notifications,
		manualRequests: make(chan struct{}, 1),
	}
	obj.status.SetStartedAt(time.Unix(0, 0))
	obj.status.SetStoppedAt(time.Unix(0, 0))
	obj.status.SetUpdatedAt(time.Unix(0, 0))
	obj.status.SetIsRunning(false)
	obj.status.SetRestarts(0)
	return obj
}
