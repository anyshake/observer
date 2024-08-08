package explorer

import (
	"context"
	"time"

	"github.com/anyshake/observer/drivers/transport"
	"github.com/anyshake/observer/utils/timesource"
	cmap "github.com/orcaman/concurrent-map/v2"
	messagebus "github.com/vardius/message-bus"
)

const EXPLORER_GENERAL_JITTER = 10 * time.Millisecond

const (
	EXPLORER_CHANNEL_CODE_Z = "Z"
	EXPLORER_CHANNEL_CODE_E = "E"
	EXPLORER_CHANNEL_CODE_N = "N"
)

type ExplorerHealth struct {
	SampleRate int
	Errors     int64
	Received   int64
	StartTime  time.Time
	UpdatedAt  time.Time // Last local system time the health information was updated
}

type ExplorerConfig struct {
	NoGeophone bool
	LegacyMode bool
	DeviceId   uint32
	Latitude   float64
	Longitude  float64
	Elevation  float64
}

type ExplorerDependency struct {
	Health       ExplorerHealth
	Config       ExplorerConfig
	FallbackTime timesource.Source
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
