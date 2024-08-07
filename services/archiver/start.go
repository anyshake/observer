package archiver

import (
	"sync"

	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/services"
	"github.com/anyshake/observer/utils/logger"
)

func (a *ArchiverService) Start(options *services.Options, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	// Get lifecycle from configuration
	serviceConfig, ok := options.Config.Services[a.GetServiceName()]
	if !ok {
		logger.GetLogger(a.GetServiceName()).Errorln("service configuration not found")
		return
	}
	if !serviceConfig.(map[string]any)["enable"].(bool) {
		logger.GetLogger(a.GetServiceName()).Infoln("service has been disabled")
		return
	}
	a.lifeCycle = int(serviceConfig.(map[string]any)["lifecycle"].(float64))
	a.cleanupCountDown = CLEANUP_COUNTDOWN
	a.insertCountDown = INSERT_COUNTDOWN
	a.databaseConn = options.Database

	// Subscribe to Explorer events
	var explorerDeps *explorer.ExplorerDependency
	err := options.Dependency.Invoke(func(deps *explorer.ExplorerDependency) error {
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
	<-options.CancelToken.Done()
	explorerDriver.Unsubscribe(explorerDeps, a.GetServiceName())
	logger.GetLogger(a.GetServiceName()).Infoln("service has been stopped")
}
