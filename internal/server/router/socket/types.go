package socket

import (
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/pkg/message"
	"github.com/anyshake/observer/pkg/ringbuf"
)

const LOG_PREFIX = "websocket_api_stream"

const HISTORY_BUFFER_SIZE = 120

type buffer struct {
	SampleRate  int
	Timestamp   int64
	ChannelData []explorer.ChannelData
}

type socket struct {
	messageBus     *message.Bus[explorer.Event]
	tokenValidator func(string) bool
	history        *ringbuf.Buffer[buffer]
}
