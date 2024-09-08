package forwarder

import (
	"github.com/alphadose/haxmap"
	"github.com/anyshake/observer/drivers/explorer"
	messagebus "github.com/vardius/message-bus"
)

type ForwarderService struct {
	messageBus    messagebus.MessageBus // An independent message bus for the socket module
	subscribers   *haxmap.Map[string, func(data *explorer.ExplorerData)]
	stationCode   string
	networkCode   string
	locationCode  string
	channelPrefix string
}
