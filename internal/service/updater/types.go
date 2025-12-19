package updater

import (
	"context"
	"sync"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/internal/upgrade"
	"github.com/anyshake/observer/pkg/timesource"
)

const ID = "service_self_updater"

const (
	UPDATE_CHECK_INTERVAL = 6 * time.Hour
	UPDATE_FETCH_TIMEOUT  = time.Hour
)

type UpdaterServiceImpl struct {
	mu     sync.Mutex
	status service.Status

	wg       sync.WaitGroup
	ctx      context.Context
	cancelFn context.CancelFunc

	actionHandler *action.Handler
	timeSource    *timesource.Source
	upgradeHelper *upgrade.Helper
}
