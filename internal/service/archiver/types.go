package archiver

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/timesource"
)

const ID = "service_archiver"

const (
	RECORDS_INSERT_INTERVAL  = 10
	RECORDS_CLEANUP_INTERVAL = 3600

	// RECORDS_FLUSH_INTERVAL bounds how long a partial batch waits in memory.
	RECORDS_FLUSH_INTERVAL = 15 * time.Second
	// RECORDS_INSERT_TIMEOUT bounds a single insert transaction.
	RECORDS_INSERT_TIMEOUT = 20 * time.Second
	// RECORDS_INSERT_RETRY_DELAY spaces retries after a failed insert.
	RECORDS_INSERT_RETRY_DELAY = 5 * time.Second
	// RECORDS_INSERT_MAX_ATTEMPTS discards a batch that keeps failing with a
	// non-transient error instead of blocking newer records forever. Timeouts
	// and busy errors are retried until the pending limit drops old records.
	RECORDS_INSERT_MAX_ATTEMPTS = 5
	// RECORDS_INSERT_BATCH_LIMIT caps one transaction when catching up.
	RECORDS_INSERT_BATCH_LIMIT = 120
	// RECORDS_PENDING_LIMIT is the in-memory backlog (one record per second)
	// kept while the database is unavailable; older records are discarded.
	RECORDS_PENDING_LIMIT = 1800
	// RECORDS_FINAL_FLUSH_TIMEOUT bounds the flush performed while stopping.
	RECORDS_FINAL_FLUSH_TIMEOUT = 5 * time.Second
)

type ArchiverServiceImpl struct {
	lifecycle service.Lifecycle
	mu        sync.Mutex
	status    service.Status

	wg        sync.WaitGroup
	parentCtx context.Context
	ctx       context.Context
	cancelFn  context.CancelFunc

	hardwareDev   hardware.IHardware
	timeSource    *timesource.Source
	actionHandler *action.Handler

	stationCode  string
	networkCode  string
	locationCode string

	rotation       int
	cleanupRunning atomic.Bool
}
