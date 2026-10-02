package tailscale

import (
	"context"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	tailscale_store "github.com/anyshake/observer/internal/service/tailscale/store"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/timesource"
	"tailscale.com/ipn"
)

func newDatabaseStateStore(actionHandler *action.Handler) (ipn.StateStore, error) {
	return tailscale_store.New(actionHandler, ID)
}

func New(webServerAddr string, actionHandler *action.Handler, timeSource *timesource.Source) *TailscaleServiceImpl {
	ctx, cancelFn := context.WithCancel(context.Background())
	obj := &TailscaleServiceImpl{
		ctx:           ctx,
		cancelFn:      cancelFn,
		actionHandler: actionHandler,
		timeSource:    timeSource,
		webServerAddr: webServerAddr,
		log:           logger.GetLogger(ID),
		newNode:       newTSNetNode,
		newStateStore: newDatabaseStateStore,
	}
	obj.status.SetStartedAt(time.Unix(0, 0))
	obj.status.SetStoppedAt(time.Unix(0, 0))
	obj.status.SetUpdatedAt(time.Unix(0, 0))
	obj.status.SetIsRunning(false)
	obj.status.SetRestarts(0)
	return obj
}
