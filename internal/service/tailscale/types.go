package tailscale

import (
	"context"
	"sync"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/timesource"
	"tailscale.com/ipn"
)

const ID = "service_tailscale"

const shutdownTimeout = 10 * time.Second

type serviceLogger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

type stateStoreFactory func(*action.Handler) (ipn.StateStore, error)

type TailscaleServiceImpl struct {
	mu       sync.Mutex
	errMu    sync.Mutex
	status   service.Status
	runError error

	ctx      context.Context
	cancelFn context.CancelFunc
	wg       sync.WaitGroup

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
