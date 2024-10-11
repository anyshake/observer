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
	defer logger.GetLogger(s.GetServiceName()).Infoln("service has been stopped")

	// To calculate duration to next 00:00:00 UTC
	calcDuration := func(currentTime time.Time) time.Duration {
		nextTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC)
		if currentTime.After(nextTime) {
			nextTime = nextTime.Add(24 * time.Hour)
		}
		return nextTime.Sub(currentTime)
	}

	currentTime := options.TimeSource.Get()
	timer := time.NewTimer(calcDuration(currentTime))
	defer timer.Stop()

	for {
		select {
		case <-options.CancelToken.Done():
			return
		case <-timer.C:
			err := options.TimeSource.Update()
			if err != nil {
				logger.GetLogger(s.GetServiceName()).Errorln(err)
				continue
			}

			// Reset timer to next update
			currentTime = options.TimeSource.Get()
			timer.Reset(calcDuration(currentTime))
			logger.GetLogger(s.GetServiceName()).Info("time source has been updated")
		}
	}
}
