package archiver

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/timesource"
)

const ID = "service_archiver"

const (
	RECORDS_INSERT_INTERVAL  = 10
	RECORDS_CLEANUP_INTERVAL = 3600
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

	rotation         int
	cleanupCountDown int
	cleanupRunning   atomic.Bool
	insertCountDown  int
	recordBuffer     []model.SeisRecord
}
