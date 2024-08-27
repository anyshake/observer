package watchdog

import (
	"sync"
	"time"

	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/services"
	"github.com/anyshake/observer/utils/logger"
)

func (s *WatchdogService) Start(options *services.Options, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	logger.GetLogger(s.GetServiceName()).Infoln("service has been started")

	var explorerDeps *explorer.ExplorerDependency
	err := options.Dependency.Invoke(func(deps *explorer.ExplorerDependency) error {
		explorerDeps = deps
		return nil
	})
	if err != nil {
		logger.GetLogger(s.GetServiceName()).Errorln(err)
		return
	}

	prevUpdatedAt := explorerDeps.Health.GetUpdatedAt()
	ticker := time.NewTicker(CHECK_INTERVAL)

	for {
		select {
		case <-options.CancelToken.Done():
			ticker.Stop()
			logger.GetLogger(s.GetServiceName()).Infoln("service has been stopped")
			return
		case <-ticker.C:
			currentUpdatedAt := explorerDeps.Health.GetUpdatedAt()
			if prevUpdatedAt == currentUpdatedAt {
				logger.GetLogger(s.GetServiceName()).Warnf("device is not responding, checking again in next %d seconds", int(CHECK_INTERVAL.Seconds()))
			}
			prevUpdatedAt = currentUpdatedAt
		}
	}
}
