package socket

import (
	"github.com/alphadose/haxmap"
	"github.com/anyshake/observer/drivers/explorer"
	messagebus "github.com/vardius/message-bus"
)

const HISTORY_BUFFER_SIZE = 180

type Socket struct {
	messageBus         messagebus.MessageBus // An independent message bus for the socket module
	subscribers        *haxmap.Map[string, func(data *explorer.ExplorerData)]
	historyBuffer      [HISTORY_BUFFER_SIZE]explorer.ExplorerData
	historyBufferIndex int
}
