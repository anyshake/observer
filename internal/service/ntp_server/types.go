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
	mu     sync.Mutex
	status service.Status

	wg       sync.WaitGroup
	ctx      context.Context
	cancelFn context.CancelFunc

	timeSource    *timesource.Source
	actionHandler *action.Handler

	listenHost string
	listenPort int
}
