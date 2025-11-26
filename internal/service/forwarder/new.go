package forwarder

import (
	"context"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/pkg/message"
	"github.com/anyshake/observer/pkg/timesource"
)

func New(hardwareDev hardware.IHardware, actionHandler *action.Handler, timeSource *timesource.Source) *ForwarderServiceImpl {
	ctx, cancelFn := context.WithCancel(context.Background())
	obj := &ForwarderServiceImpl{
		hardwareDev:        hardwareDev,
		ctx:                ctx,
		cancelFn:           cancelFn,
		timeSource:         timeSource,
		actionHandler:      actionHandler,
		messageBus:         message.NewBus[explorer.EventHandler](ID, 65535),
		messageBusRealtime: message.NewBus[explorer.EventHandler](ID, 65535),
	}
	obj.status.SetStartedAt(time.Unix(0, 0))
	obj.status.SetStoppedAt(time.Unix(0, 0))
	obj.status.SetUpdatedAt(time.Unix(0, 0))
	obj.status.SetIsRunning(false)
	obj.status.SetRestarts(0)
	return obj
}
