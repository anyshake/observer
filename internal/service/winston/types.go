package winston

import (
	"context"
	"sync"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/ringbuf"
	"github.com/anyshake/observer/pkg/timesource"
)

const ID = "service_winston"

type winstonRingBuffer struct {
	timestamp   time.Time
	sampleRate  int
	channelData []explorer.ChannelData
}

type WinstonServiceImpl struct {
	lifecycle service.Lifecycle
	mu        sync.Mutex
	status    service.Status

	wg        sync.WaitGroup
	parentCtx context.Context
	ctx       context.Context
	cancelFn  context.CancelFunc

	hardwareDev   hardware.IHardware
	timeSource    *timesource.Source
	actionHandler *action.Handler

	stationCode  string
	networkCode  string
	locationCode string

	listenHost string
	listenPort int
	bufferSize int

	// Holds the most recent waveform records for low-latency queries.
	ringBuffer *ringbuf.Buffer[winstonRingBuffer]
}
