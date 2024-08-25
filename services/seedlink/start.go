package seedlink

import (
	"sync"

	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/services"
	"github.com/anyshake/observer/utils/logger"
	"github.com/bclswl0827/slgo"
	cmap "github.com/orcaman/concurrent-map/v2"
	messagebus "github.com/vardius/message-bus"
)

func (s *SeedLinkService) Start(options *services.Options, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	serviceConfig, ok := options.Config.Services[s.GetServiceName()]
	if !ok {
		logger.GetLogger(s.GetServiceName()).Errorln("service configuration not found")
		return
	}
	if !serviceConfig.(map[string]any)["enable"].(bool) {
		logger.GetLogger(s.GetServiceName()).Infoln("service has been disabled")
		return
	}
	serverHost := serviceConfig.(map[string]any)["host"].(string)
	serverPort := int(serviceConfig.(map[string]any)["port"].(float64))
	currentTime, _ := options.TimeSource.Get()
	messageBus := messagebus.New(65535)

	// Subscribe to Explorer events
	var explorerDeps *explorer.ExplorerDependency
	err := options.Dependency.Invoke(func(deps *explorer.ExplorerDependency) error {
		explorerDeps = deps
		return nil
	})
	if err != nil {
		logger.GetLogger(s.GetServiceName()).Errorln(err)
		return
	}
	explorerDriver := explorer.ExplorerDriver(&explorer.ExplorerDriverImpl{})
	explorerDriver.Subscribe(
		explorerDeps,
		s.GetServiceName(),
		func(data *explorer.ExplorerData) {
			if s.prevSampleRate == 0 {
				s.prevSampleRate = data.SampleRate
			}
			if s.prevSampleRate == data.SampleRate {
				messageBus.Publish(s.GetServiceName(), data)
			} else {
				logger.GetLogger(s.GetServiceName()).Warnf("sample rate is not the same, expected %d, got %d", s.prevSampleRate, data.SampleRate)
			}
			s.prevSampleRate = data.SampleRate
		},
	)

	// Start SeedLink server
	server := slgo.New(
		&provider{
			timeSource:    options.TimeSource,
			database:      options.Database,
			startTime:     currentTime,
			stationCode:   options.Config.Stream.Station,
			networkCode:   options.Config.Stream.Network,
			locationCode:  options.Config.Stream.Location,
			channelPrefix: options.Config.Stream.Channel,
		},
		&consumer{
			channelPrefix: options.Config.Stream.Channel,
			serviceName:   s.GetServiceName(),
			messageBus:    messageBus,
			subscribers:   cmap.New[explorer.ExplorerEventHandler](),
		},
		&hooks{
			serviceName: s.GetServiceName(),
		},
	)
	go server.Start(serverHost, serverPort)
	logger.GetLogger(s.GetServiceName()).Infof("seedlink is listening on %s:%d", serverHost, serverPort)

	logger.GetLogger(s.GetServiceName()).Infoln("service has been started")
	defer logger.GetLogger(s.GetServiceName()).Infoln("service has been stopped")

	<-options.CancelToken.Done()
	explorerDriver.Unsubscribe(explorerDeps, s.GetServiceName())
}
