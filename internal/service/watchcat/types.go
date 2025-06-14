package watchcat

import (
	"context"
	"sync"
	"time"

	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/timesource"
)

const ID = "service_watchcat"

const (
	WATCHCAT_CHECK_INTERVAL = 30 * time.Second
)

type WatchCatServiceImpl struct {
	mu     sync.Mutex
	status service.Status

	wg       sync.WaitGroup
	ctx      context.Context
	cancelFn context.CancelFunc

	hardwareDev hardware.IHardware
	timeSource  *timesource.Source
}
