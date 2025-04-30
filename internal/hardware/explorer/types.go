package explorer

import (
	"sync"
	"time"
)

const (
	EXPLORER_STREAM_TOPIC = "/explorer/stream"
	ALLOWED_JITTER_MS     = 10
)

type DeviceStatus struct {
	mu sync.Mutex

	startedAt time.Time
	updatedAt time.Time
	frames    int64
	errors    int64
	messages  int64
}

type DeviceConfig struct {
	mu sync.Mutex

	packetInterval time.Duration
	channelCodes   []string
	sampleRate     int
	gnssEnabled    bool
	model          string
	protocol       string
}

type DeviceVariable struct {
	mu sync.Mutex

	deviceId    *uint32
	latitude    *float64
	longitude   *float64
	elevation   *float64
	temperature *float64
}

type ChannelData struct {
	ChannelCode string
	ChannelId   int
	ByteSize    int
	DataType    string
	Data        []int32
}

type EventHandler = func(time.Time, *DeviceConfig, *DeviceVariable, []ChannelData)
