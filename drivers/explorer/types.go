package explorer

import (
	"context"
	"sync"
	"time"

	"github.com/alphadose/haxmap"
	"github.com/anyshake/observer/drivers/transport"
	"github.com/anyshake/observer/utils/timesource"
	"github.com/sirupsen/logrus"
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
	mutex             sync.RWMutex
	legacyMode        bool
	deviceInfo        uint32
	gnssLatitude      float64
	fallbackLatitude  float64
	gnssLongitude     float64
	fallbackLongitude float64
	gnssElevation     float64
	fallbackElevation float64
}

type ExplorerDependency struct {
	FallbackTime *timesource.Source
	Health       *ExplorerHealth
	Config       *ExplorerConfig
	CancelToken  context.Context
	Transport    transport.TransportDriver
	messageBus   messagebus.MessageBus
	subscribers  *haxmap.Map[string, func(data *ExplorerData)]
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
	Init(deps *ExplorerDependency, logger *logrus.Entry) error
	Subscribe(deps *ExplorerDependency, clientId string, handler ExplorerEventHandler) error
	Unsubscribe(deps *ExplorerDependency, clientId string) error
}
