package explorer

import (
	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/startups"
	"github.com/anyshake/observer/utils/logger"
	"go.uber.org/dig"
)

func (t *ExplorerStartupTask) Execute(depsContainer *dig.Container, options *startups.Options) error {
	var explorerDeps *explorer.ExplorerDependency
	err := depsContainer.Invoke(func(deps *explorer.ExplorerDependency) error {
		explorerDeps = deps
		return nil
	})
	if err != nil {
		return err
	}
	explorerDriver := explorer.ExplorerDriver(&explorer.ExplorerDriverImpl{})

	logger.GetLogger(t.GetTaskName()).Infoln("device is being initialized, please wait")
	err = explorerDriver.Init(explorerDeps, logger.GetLogger("explorer_driver"))
	if err != nil {
		return err
	}

	logger.GetLogger(t.GetTaskName()).Infoln("device has been initialized successfully")
	if !explorerDeps.Config.GetLegacyMode() {
		logger.GetLogger(t.GetTaskName()).Infof("handshake successful, device ID: %08X", explorerDeps.Config.GetDeviceId())
	} else {
		logger.GetLogger(t.GetTaskName()).Warnln("device is in legacy mode, this is for backward compatibility only")
	}
	return nil
}
