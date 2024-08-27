package explorer

import (
	"context"
	"sync"
	"time"

	"github.com/anyshake/observer/drivers/transport"
	"github.com/anyshake/observer/utils/timesource"
	cmap "github.com/orcaman/concurrent-map/v2"
	messagebus "github.com/vardius/message-bus"
)

const EXPLORER_ALLOWED_JITTER_MS = 5

const (
	EXPLORER_CHANNEL_CODE_Z = "Z"
	EXPLORER_CHANNEL_CODE_E = "E"
	EXPLORER_CHANNEL_CODE_N = "N"
)

type ExplorerHealth struct {
	mutex      sync.RWMutex
	sampleRate int
	errors     int64
	received   int64
	startTime  time.Time
	updatedAt  time.Time // Last local system time the health information was updated
}

type ExplorerConfig struct {
	mutex      sync.RWMutex
	legacyMode bool
	deviceId   uint32
	latitude   float64
	longitude  float64
	elevation  float64
}

type ExplorerDependency struct {
	FallbackTime *timesource.Source
	Health       *ExplorerHealth
	Config       *ExplorerConfig
	CancelToken  context.Context
	Transport    transport.TransportDriver
	messageBus   messagebus.MessageBus
	subscribers  cmap.ConcurrentMap[string, ExplorerEventHandler]
}

type ExplorerData struct {
	SampleRate int     `json:"sample_rate"`
	Timestamp  int64   `json:"timestamp"`
	Z_Axis     []int32 `json:"z_axis"`
	E_Axis     []int32 `json:"e_axis"`
	N_Axis     []int32 `json:"n_axis"`
}

type ExplorerEventHandler = func(data *ExplorerData)

type ExplorerDriver interface {
	readerDaemon(deps *ExplorerDependency)
	IsAvailable(deps *ExplorerDependency) bool
	Init(deps *ExplorerDependency) error
	Subscribe(deps *ExplorerDependency, clientId string, handler ExplorerEventHandler) error
	Unsubscribe(deps *ExplorerDependency, clientId string) error
}
