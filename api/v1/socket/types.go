package socket

import (
	"github.com/anyshake/observer/drivers/explorer"
	cmap "github.com/orcaman/concurrent-map/v2"
	messagebus "github.com/vardius/message-bus"
)

const EXPLORER_BUFFER_SIZE = 180

type Socket struct {
	messageBus         messagebus.MessageBus // An independent message bus for the socket module
	subscribers        cmap.ConcurrentMap[string, explorer.ExplorerEventHandler]
	historyBuffer      [EXPLORER_BUFFER_SIZE]explorer.ExplorerData
	historyBufferIndex int
}
