package miniseed

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/services"
	"github.com/anyshake/observer/utils/logger"
	"github.com/bclswl0827/mseedio"
)

func (m *MiniSeedService) Start(options *services.Options, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	// Get lifecycle from configuration
	serviceConfig, ok := options.Config.Services[m.GetServiceName()]
	if !ok {
		logger.GetLogger(m.GetServiceName()).Errorln("service configuration not found")
		return
	}
	if !serviceConfig.(map[string]any)["enable"].(bool) {
		logger.GetLogger(m.GetServiceName()).Infoln("service has been disabled")
		return
	}
	m.lifeCycle = int(serviceConfig.(map[string]any)["lifecycle"].(float64))
	m.basePath = serviceConfig.(map[string]any)["path"].(string)
	m.stationCode = options.Config.Stream.Station
	m.networkCode = options.Config.Stream.Network
	m.locationCode = options.Config.Stream.Location
	m.channelPrefix = options.Config.Stream.Channel
	m.miniseedSequence = map[string]int{
		explorer.EXPLORER_CHANNEL_CODE_Z: 0,
		explorer.EXPLORER_CHANNEL_CODE_E: 0,
		explorer.EXPLORER_CHANNEL_CODE_N: 0,
	}
	m.cleanUpCountDown = MINISEED_CLEANUP_INTERVAL
	m.writeBufferCountDown = MINISEED_WRITE_INTERVAL

	// Get sequence number if file exists
	currentTime, _ := options.TimeSource.GetTime()
	for _, channelCode := range []string{
		explorer.EXPLORER_CHANNEL_CODE_Z,
		explorer.EXPLORER_CHANNEL_CODE_E,
		explorer.EXPLORER_CHANNEL_CODE_N,
	} {
		channelName := fmt.Sprintf("%s%s", m.channelPrefix, channelCode)
		filePath := m.getFilePath(
			m.basePath,
			m.stationCode,
			m.networkCode,
			m.locationCode,
			channelName,
			currentTime,
		)
		_, err := os.Stat(filePath)
		if err == nil {
			logger.GetLogger(m.GetServiceName()).Infof("starting %s from last record", channelName)

			// Get last sequence number from file
			var ms mseedio.MiniSeedData
			err := ms.Read(filePath)
			if err != nil {
				logger.GetLogger(m.GetServiceName()).Errorln(err)
				continue
			}
			if len(ms.Series) > 0 {
				lastRecord := ms.Series[len(ms.Series)-1]
				lastSeqNum := lastRecord.FixedSection.SequenceNumber
				n, err := strconv.Atoi(lastSeqNum)
				if err != nil {
					logger.GetLogger(m.GetServiceName()).Errorln(err)
					continue
				}

				// Set current sequence number
				m.miniseedSequence[channelCode] = n
			}
		}
	}

	// Subscribe to Explorer events
	var explorerDeps *explorer.ExplorerDependency
	err := options.Dependency.Invoke(func(deps *explorer.ExplorerDependency) error {
		explorerDeps = deps
		return nil
	})
	if err != nil {
		logger.GetLogger(m.GetServiceName()).Errorln(err)
		return
	}
	explorerDriver := explorer.ExplorerDriver(&explorer.ExplorerDriverImpl{})
	explorerDriver.Subscribe(explorerDeps, m.GetServiceName(), m.handleExplorerEvent)

	logger.GetLogger(m.GetServiceName()).Infoln("service has been started")
	<-options.CancelToken.Done()
	explorerDriver.Unsubscribe(explorerDeps, m.GetServiceName())
	logger.GetLogger(m.GetServiceName()).Infoln("service has been stopped")
}
