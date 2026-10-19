package ntp_server

import (
	"context"
	"sync"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/timesource"
)

const ID = "service_ntp_server"

type NtpServerServiceImpl struct {
	lifecycle service.Lifecycle
	mu        sync.Mutex
	status    service.Status

	wg        sync.WaitGroup
	parentCtx context.Context
	ctx       context.Context
	cancelFn  context.CancelFunc

	timeSource    *timesource.Source
	actionHandler *action.Handler

	listenHost   string
	listenPort   int
	delayControl int64
}
