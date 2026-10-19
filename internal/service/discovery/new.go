package mdns_discovery

import (
	"context"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/timesource"
)

func New(parentCtx context.Context, webServerAddr string, serviceMap map[string]service.IService, actionHandler *action.Handler, timeSource *timesource.Source) *DiscoveryServiceImpl {
	ctx, cancelFn := context.WithCancel(parentCtx)
	obj := &DiscoveryServiceImpl{
		parentCtx:     parentCtx,
		ctx:           ctx,
		cancelFn:      cancelFn,
		timeSource:    timeSource,
		actionHandler: actionHandler,
		webServerAddr: webServerAddr,
		serviceMap:    serviceMap,
	}
	obj.status.SetStartedAt(time.Unix(0, 0))
	obj.status.SetStoppedAt(time.Unix(0, 0))
	obj.status.SetUpdatedAt(time.Unix(0, 0))
	obj.status.SetIsRunning(false)
	obj.status.SetRestarts(0)
	return obj
}
