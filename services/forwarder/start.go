package forwarder

import (
	"fmt"
	"net"
	"sync"

	"github.com/alphadose/haxmap"
	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/services"
	"github.com/anyshake/observer/utils/logger"
	messagebus "github.com/vardius/message-bus"
)

func (a *ForwarderService) Start(options *services.Options, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	enabled, err := options.Config.Services.GetValue(a.GetServiceName(), "enable", "bool")
	if err != nil {
		logger.GetLogger(a.GetServiceName()).Errorln(err)
		return
	}
	if !enabled.(bool) {
		logger.GetLogger(a.GetServiceName()).Infoln("service has been disabled")
		return
	}

	serverHost, err := options.Config.Services.GetValue(a.GetServiceName(), "host", "string")
	if err != nil {
		logger.GetLogger(a.GetServiceName()).Errorln(err)
		return
	}
	serverPort, err := options.Config.Services.GetValue(a.GetServiceName(), "port", "int")
	if err != nil {
		logger.GetLogger(a.GetServiceName()).Errorln(err)
		return
	}

	a.stationCode = options.Config.Stream.Station
	a.networkCode = options.Config.Stream.Network
	a.locationCode = options.Config.Stream.Location
	a.channelPrefix = options.Config.Stream.Channel

	a.subscribers = haxmap.New[string, explorer.ExplorerEventHandler]()
	a.messageBus = messagebus.New(65535)

	// Forward events to internal message bus
	var explorerDeps *explorer.ExplorerDependency
	err = options.Dependency.Invoke(func(deps *explorer.ExplorerDependency) error {
		explorerDeps = deps
		return nil
	})
	if err != nil {
		logger.GetLogger(a.GetServiceName()).Errorln(err)
		return
	}
	explorerDriver := explorer.ExplorerDriver(&explorer.ExplorerDriverImpl{})
	explorerDriver.Subscribe(
		explorerDeps,
		a.GetServiceName(),
		func(data *explorer.ExplorerData) { a.messageBus.Publish(a.GetServiceName(), data) },
	)

	// Create TCP server to forward events
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serverHost.(string), serverPort.(int)))
	if err != nil {
		logger.GetLogger(a.GetServiceName()).Errorln(err)
		return
	}
	defer listener.Close()
	logger.GetLogger(a.GetServiceName()).Infof("forwarder is listening on %s:%d", serverHost, serverPort)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go a.handleConnection(conn)
		}
	}()

	logger.GetLogger(a.GetServiceName()).Infoln("service has been started")
	defer logger.GetLogger(a.GetServiceName()).Infoln("service has been stopped")

	<-options.CancelToken.Done()
	explorerDriver.Unsubscribe(explorerDeps, a.GetServiceName())
}
