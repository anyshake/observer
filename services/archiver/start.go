package archiver

import (
	"sync"

	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/services"
	"github.com/anyshake/observer/utils/logger"
)

func (a *ArchiverService) Start(options *services.Options, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	enabled, err := options.Config.Services.GetValue(a.GetServiceName(), "enable", "bool")
	if err != nil {
		logger.GetLogger(a.GetServiceName()).Errorln(err)
		return
	}
	if !enabled.(bool) {
		logger.GetLogger(a.GetServiceName()).Infoln("service has been disabled")
		return
	}
	lifecycle, err := options.Config.Services.GetValue(a.GetServiceName(), "lifecycle", "int")
	if err != nil {
		logger.GetLogger(a.GetServiceName()).Errorln(err)
		return
	}

	a.lifeCycle = lifecycle.(int)
	a.cleanupCountDown = RECORDS_CLEANUP_INTERVAL
	a.insertCountDown = RECORDS_INSERT_INTERVAL
	a.databaseConn = options.Database

	// Subscribe to Explorer events
	var explorerDeps *explorer.ExplorerDependency
	err = options.Dependency.Invoke(func(deps *explorer.ExplorerDependency) error {
		explorerDeps = deps
		return nil
	})
	if err != nil {
		logger.GetLogger(a.GetServiceName()).Errorln(err)
		return
	}
	explorerDriver := explorer.ExplorerDriver(&explorer.ExplorerDriverImpl{})
	explorerDriver.Subscribe(explorerDeps, a.GetServiceName(), a.handleExplorerEvent)

	logger.GetLogger(a.GetServiceName()).Infoln("service has been started")
	defer logger.GetLogger(a.GetServiceName()).Infoln("service has been stopped")

	<-options.CancelToken.Done()
	explorerDriver.Unsubscribe(explorerDeps, a.GetServiceName())
}
