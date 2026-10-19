package tailscale

import (
	"context"
	"sync"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/timesource"
	"tailscale.com/ipn"
)

const ID = "service_tailscale"

type serviceLogger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

type stateStoreFactory func(context.Context, *action.Handler) (ipn.StateStore, error)

type TailscaleServiceImpl struct {
	lifecycle service.Lifecycle
	mu        sync.Mutex
	status    service.Status

	parentCtx context.Context
	ctx       context.Context
	cancelFn  context.CancelFunc
	wg        sync.WaitGroup

	timeSource    *timesource.Source
	actionHandler *action.Handler
	webServerAddr string

	stateStore    ipn.StateStore
	log           serviceLogger
	newNode       nodeFactory
	newStateStore stateStoreFactory

	hostname     string
	controlURL   string
	authKey      string
	forwardPorts []string
}
