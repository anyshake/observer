package mdns_discovery

import (
	"context"
	"sync"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/grandcat/zeroconf"
)

const ID = "service_mdns_discovery"

type DiscoveryServiceImpl struct {
	mu     sync.Mutex
	status service.Status

	wg       sync.WaitGroup
	ctx      context.Context
	cancelFn context.CancelFunc

	timeSource      *timesource.Source
	actionHandler   *action.Handler
	localServerAddr string

	instanceName string
	server       *zeroconf.Server
}
