package frp_client

import (
	"context"
	"strings"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/fatedier/frp/pkg/util/log"
	loglib "github.com/fatedier/golib/log"
)

type logWriter struct{}

func (w *logWriter) Write(p []byte) (n int, err error) {
	logText := string(p)
	logText = strings.TrimSuffix(logText, "\n")

	if len(logText) < 28 {
		return len(p), nil
	}

	switch logText[25:26] {
	case "I":
		logger.GetLogger(ID).Infof("FRP client returned info: %s", logText[28:])
	case "W":
		logger.GetLogger(ID).Warnf("FRP client returned warning: %s", logText[28:])
	case "E":
		logger.GetLogger(ID).Errorf("FRP client returned error: %s", logText[28:])
	}

	return len(p), nil
}

func New(localServerAddr string, actionHandler *action.Handler, timeSource *timesource.Source) *FrpClientServiceImpl {
	log.Logger = loglib.New(
		loglib.WithLevel(log.InfoLevel),
		loglib.WithOutput(loglib.NewConsoleWriter(loglib.ConsoleConfig{Colorful: false}, &logWriter{})),
	)

	ctx, cancelFn := context.WithCancel(context.Background())
	obj := &FrpClientServiceImpl{
		ctx:             ctx,
		cancelFn:        cancelFn,
		timeSource:      timeSource,
		actionHandler:   actionHandler,
		localServerAddr: localServerAddr,
	}
	obj.status.SetStartedAt(time.Unix(0, 0))
	obj.status.SetStoppedAt(time.Unix(0, 0))
	obj.status.SetUpdatedAt(time.Unix(0, 0))
	obj.status.SetIsRunning(false)
	obj.status.SetRestarts(0)
	return obj
}
