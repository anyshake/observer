package explorer

import (
	"sync"
	"time"

	"github.com/anyshake/observer/pkg/message"
)

const (
	EXPLORER_STREAM_TOPIC          = "/explorer/stream"          // 1 message per second
	EXPLORER_REALTIME_STREAM_TOPIC = "/explorer/stream/realtime" // 1 message per packet
)

const (
	NTP_RESYNC_INTERVAL      = 2 * time.Minute
	NTP_MEASUREMENT_ATTEMPTS = 5
	NTP_PPM_MEASURE_WINDOW   = time.Hour
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

type Event struct {
	Timestamp     time.Time
	SampleRate    int
	GNSSAvailable bool
	ChannelData   []ChannelData
}

type EventHandler = func(Event)

func NewEvent(timestamp time.Time, deviceConfig *DeviceConfig, channelData []ChannelData) Event {
	clonedChannelData := make([]ChannelData, len(channelData))
	for i := range channelData {
		clonedChannelData[i] = channelData[i]
		clonedChannelData[i].Data = append([]int32(nil), channelData[i].Data...)
	}

	return Event{
		Timestamp:     timestamp,
		SampleRate:    deviceConfig.GetSampleRate(),
		GNSSAvailable: deviceConfig.GetGnssAvailability(),
		ChannelData:   clonedChannelData,
	}
}

func normalStreamSubscriptionOptions(onError func(error)) message.SubscriptionOptions {
	return message.SubscriptionOptions{
		BufferSize:     32,
		Overflow:       message.OverflowBlock,
		EnqueueTimeout: 250 * time.Millisecond,
		OnError:        onError,
	}
}

func realtimeStreamSubscriptionOptions(onError func(error)) message.SubscriptionOptions {
	return message.SubscriptionOptions{
		BufferSize: 8,
		Overflow:   message.OverflowDropOldest,
		OnError:    onError,
	}
}
