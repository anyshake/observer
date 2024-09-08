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

	enabled, err := options.Config.Services.GetValue(m.GetServiceName(), "enable", "bool")
	if err != nil {
		logger.GetLogger(m.GetServiceName()).Errorln(err)
		return
	}
	if !enabled.(bool) {
		logger.GetLogger(m.GetServiceName()).Infoln("service has been disabled")
		return
	}

	lifecycle, err := options.Config.Services.GetValue(m.GetServiceName(), "lifecycle", "int")
	if err != nil {
		logger.GetLogger(m.GetServiceName()).Errorln(err)
		return
	}
	m.lifeCycle = lifecycle.(int)

	basePath, err := options.Config.Services.GetValue(m.GetServiceName(), "path", "string")
	if err != nil {
		logger.GetLogger(m.GetServiceName()).Errorln(err)
		return
	}
	m.basePath = basePath.(string)

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

	// Set write interval to 1 if legacy mode is enabled
	// This is a simple solution to sample rate jittering
	// However, this will increase the disk I/O and file size
	if options.Config.Explorer.Legacy {
		m.writeBufferInterval = 1
	} else {
		m.writeBufferInterval = MINISEED_WRITE_INTERVAL
	}
	m.writeBufferCountDown = m.writeBufferInterval
	m.miniseedBuffer = make([]explorer.ExplorerData, m.writeBufferInterval)

	// Get sequence number if file exists
	currentTime := options.TimeSource.Get()
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
			logger.GetLogger(m.GetServiceName()).Infof("reading existing record from %s", filePath)

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
				logger.GetLogger(m.GetServiceName()).Infof("starting %s from last record, sequence %d", channelName, n)
			}
		}
	}

	// Subscribe to Explorer events
	var explorerDeps *explorer.ExplorerDependency
	err = options.Dependency.Invoke(func(deps *explorer.ExplorerDependency) error {
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
	defer logger.GetLogger(m.GetServiceName()).Infoln("service has been stopped")

	<-options.CancelToken.Done()
	explorerDriver.Unsubscribe(explorerDeps, m.GetServiceName())
}
