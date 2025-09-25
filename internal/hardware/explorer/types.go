package explorer

import (
	"sync"
	"time"
)

const (
	EXPLORER_STREAM_TOPIC          = "/explorer/stream"          // 1 message per second
	EXPLORER_REALTIME_STREAM_TOPIC = "/explorer/stream/realtime" // 1 message per packet
)

const (
	NTP_RESYNC_INTERVAL      = 5 * time.Minute
	NTP_MEASUREMENT_ATTEMPTS = 5
)

const (
	STABLE_CHECK_SAMPLES   = 10
	ALLOWED_JITTER_MS_GNSS = 10
	ALLOWED_JITTER_MS_NTP  = 20
)

type ExplorerOptions struct {
	Endpoint    string
	Protocol    string
	Model       string
	Latitude    float64
	Longitude   float64
	Elevation   float64
	ReadTimeout int
}

type NtpOptions struct {
	Pool        []string
	Retry       int
	ReadTimeout int
}

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
