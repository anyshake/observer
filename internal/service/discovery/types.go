package mdns_discovery

import (
	"context"
	"sync"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/grandcat/zeroconf"
)

const (
	ID                      = "service_mdns_discovery"
	DISCOVERY_SCAN_INTERVAL = 10 * time.Second
)

type mdnsRegistration struct {
	port   int
	server *zeroconf.Server
}

type DiscoveryServiceImpl struct {
	lifecycle service.Lifecycle
	mu        sync.Mutex
	status    service.Status

	wg        sync.WaitGroup
	parentCtx context.Context
	ctx       context.Context
	cancelFn  context.CancelFunc

	timeSource    *timesource.Source
	actionHandler *action.Handler
	webServerAddr string
	serviceMap    map[string]service.IService

	instanceName      string
	registerSeedLink  bool
	registerWinston   bool
	registerForwarder bool

	httpReg      *mdnsRegistration
	seedlinkReg  *mdnsRegistration
	winstonReg   *mdnsRegistration
	forwarderReg *mdnsRegistration
}
