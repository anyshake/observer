package timesync

import (
	"context"
	"sync"
	"time"

	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/timesource"
)

const ID = "service_timesync"

const TIMESOURCE_REFRESH_INTERVAL = 12 * time.Hour

type TimeSyncServiceImpl struct {
	mu     sync.Mutex
	status service.Status

	wg       sync.WaitGroup
	ctx      context.Context
	cancelFn context.CancelFunc

	timeSource *timesource.Source
}
