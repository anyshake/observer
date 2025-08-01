package ntp_server

import (
	"context"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/pkg/timesource"
)

func New(actionHandler *action.Handler, timeSource *timesource.Source) *NtpServerServiceImpl {
	ctx, cancelFn := context.WithCancel(context.Background())
	obj := &NtpServerServiceImpl{
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
