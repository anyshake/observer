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

	var explorerDeps *explorer.ExplorerDependency
	err := options.Dependency.Invoke(func(deps *explorer.ExplorerDependency) error {
		explorerDeps = deps
		return nil
	})
	if err != nil {
		logger.GetLogger(s.GetServiceName()).Errorln(err)
		return
	}

	logger.GetLogger(s.GetServiceName()).Infoln("service has been started")
	defer logger.GetLogger(s.GetServiceName()).Infoln("service has been stopped")

	prevUpdatedAt := explorerDeps.Health.GetUpdatedAt()
	ticker := time.NewTicker(CHECK_INTERVAL)
	defer ticker.Stop()

	for {
		select {
		case <-options.CancelToken.Done():
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
