package socket

import (
	"sync"

	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/pkg/message"
)

const LOG_PREFIX = "websocket_api_stream"

const HISTORY_BUFFER_SIZE = 120

type buffer struct {
	SampleRate  int
	Timestamp   int64
	ChannelData []explorer.ChannelData
}

type socket struct {
	mu            sync.Mutex
	messageBus    message.Bus[explorer.EventHandler]
	historyBuffer []buffer
}
