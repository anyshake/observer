package miniseed

import (
	"sync"

	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/services"
	"github.com/anyshake/observer/utils/logger"
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

	noCompress, err := options.Config.Services.GetValue(m.GetServiceName(), "nocompress", "bool")
	if err != nil {
		noCompress = false
	}
	m.noCompress = noCompress.(bool)

	m.stationCode = options.Config.Stream.Station
	m.networkCode = options.Config.Stream.Network
	m.locationCode = options.Config.Stream.Location
	m.channelPrefix = options.Config.Stream.Channel
	m.miniseedSequence = map[string]int{
		explorer.EXPLORER_CHANNEL_CODE_Z: 1,
		explorer.EXPLORER_CHANNEL_CODE_E: 1,
		explorer.EXPLORER_CHANNEL_CODE_N: 1,
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

	// Try to get sequence numbers
	seqNums, err := m.readSequence()
	if err == nil {
		for channelCode, seqNum := range seqNums {
			m.miniseedSequence[channelCode] = seqNum
			logger.GetLogger(m.GetServiceName()).Infof(
				"starting %s%s from last record, sequence %06d",
				m.channelPrefix,
				channelCode,
				seqNum,
			)
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
