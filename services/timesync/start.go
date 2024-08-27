package timesync

import (
	"sync"
	"time"

	"github.com/anyshake/observer/services"
	"github.com/anyshake/observer/utils/logger"
)

func (s *TimeSyncService) Start(options *services.Options, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	logger.GetLogger(s.GetServiceName()).Infoln("service has been started")

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-options.CancelToken.Done():
			ticker.Stop()
			logger.GetLogger(s.GetServiceName()).Infoln("service has been stopped")
			return
		case <-ticker.C:
			currentTime, err := options.TimeSource.Get()
			if err != nil {
				logger.GetLogger(s.GetServiceName()).Errorln(err)
				continue
			}

			// Update time source at 00:00:00 UTC every day
			if currentTime.Unix()%86400 == 0 {
				if err = options.TimeSource.Update(); err != nil {
					logger.GetLogger(s.GetServiceName()).Errorln(err)
				} else {
					logger.GetLogger(s.GetServiceName()).Info("time source has been updated")
				}
			}
		}
	}
}
