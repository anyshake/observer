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
	historyMu      sync.RWMutex
	messageBus     *message.Bus[explorer.Event]
	tokenValidator func(string) bool
	historyBuffer  []buffer
	historyPos     int
	historyLen     int
}
