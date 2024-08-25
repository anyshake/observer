package forwarder

import (
	"github.com/anyshake/observer/drivers/explorer"
	cmap "github.com/orcaman/concurrent-map/v2"
	messagebus "github.com/vardius/message-bus"
)

type ForwarderService struct {
	messageBus    messagebus.MessageBus // An independent message bus for the socket module
	subscribers   cmap.ConcurrentMap[string, explorer.ExplorerEventHandler]
	stationCode   string
	networkCode   string
	locationCode  string
	channelPrefix string
}
