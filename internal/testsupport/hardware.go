package testsupport

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/pkg/metadata"
)

// Hardware is an in-memory hardware.IHardware. It never opens a device.
type Hardware struct {
	mu sync.Mutex

	Config explorer.DeviceConfig
	Status explorer.DeviceStatus

	OpenErr              error
	CloseErr             error
	FlushErr             error
	SubscribeErr         error
	SubscribeRealtimeErr error
	CoordErr             error
	TempErr              error
	MetaErr              error

	Latitude    float64
	Longitude   float64
	Elevation   float64
	Temperature float64
	DeviceID    string
	Flushes     int

	stream   explorer.EventHandler
	realtime explorer.EventHandler

	subscribed         chan struct{}
	subscribedRealtime chan struct{}
}

func NewHardware() *Hardware {
	hw := &Hardware{
		DeviceID:           "device-1",
		Latitude:           35.0,
		Longitude:          139.0,
		Elevation:          40,
		Temperature:        21.5,
		subscribed:         make(chan struct{}, 1),
		subscribedRealtime: make(chan struct{}, 1),
	}
	hw.Config.SetSampleRate(100)
	hw.Config.SetChannelCodes([]string{"EHZ", "EHN", "EHE"})
	hw.Config.SetModel("test-model")
	hw.Config.SetProtocol("v1")
	hw.Config.SetGnssAvailability(false)
	hw.Status.SetUpdatedAt(time.Unix(0, 0))
	return hw
}

func (h *Hardware) Open(ctx context.Context) (context.Context, context.CancelFunc, error) {
	if h.OpenErr != nil {
		return nil, nil, h.OpenErr
	}
	if ctx == nil {
		ctx = context.Background()
	}
	next, cancel := context.WithCancel(ctx)
	return next, cancel, nil
}

func (h *Hardware) Close() error { return h.CloseErr }

func (h *Hardware) Flush() error {
	h.mu.Lock()
	h.Flushes++
	h.mu.Unlock()
	return h.FlushErr
}

func (h *Hardware) Subscribe(_ string, handler explorer.EventHandler) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.SubscribeErr != nil {
		return h.SubscribeErr
	}
	return h.remember(&h.stream, h.subscribed, handler)
}

func (h *Hardware) remember(dst *explorer.EventHandler, ready chan struct{}, handler explorer.EventHandler) error {
	*dst = handler
	select {
	case ready <- struct{}{}:
	default:
	}
	return nil
}

func (h *Hardware) Unsubscribe(string) error {
	h.mu.Lock()
	h.stream = nil
	h.mu.Unlock()
	return nil
}

func (h *Hardware) SubscribeRealtime(_ string, handler explorer.EventHandler) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.SubscribeRealtimeErr != nil {
		return h.SubscribeRealtimeErr
	}
	return h.remember(&h.realtime, h.subscribedRealtime, handler)
}

func (h *Hardware) UnsubscribeRealtime(string) error {
	h.mu.Lock()
	h.realtime = nil
	h.mu.Unlock()
	return nil
}

func (h *Hardware) GetConfig() explorer.DeviceConfig { return h.Config }

func (h *Hardware) GetStatus() explorer.DeviceStatus { return h.Status }

func (h *Hardware) GetCoordinates(bool) (float64, float64, float64, error) {
	if h.CoordErr != nil {
		return 0, 0, 0, h.CoordErr
	}
	return h.Latitude, h.Longitude, h.Elevation, nil
}

func (h *Hardware) GetTemperature() (float64, error) {
	if h.TempErr != nil {
		return 0, h.TempErr
	}
	return h.Temperature, nil
}

func (h *Hardware) GetDeviceId() string { return h.DeviceID }

func (h *Hardware) GetMetadata(string, string, string, string, string, string, string, bool) (*metadata.Render, error) {
	if h.MetaErr != nil {
		return nil, h.MetaErr
	}
	return &metadata.Render{}, nil
}

func (h *Hardware) Emit(event explorer.Event) {
	h.mu.Lock()
	handler := h.stream
	h.mu.Unlock()
	if handler != nil {
		handler(event)
	}
}

func (h *Hardware) EmitRealtime(event explorer.Event) {
	h.mu.Lock()
	handler := h.realtime
	h.mu.Unlock()
	if handler != nil {
		handler(event)
	}
}

func (h *Hardware) WaitSubscribed(t *testing.T) {
	t.Helper()
	h.wait(t, h.subscribed)
}

func (h *Hardware) WaitSubscribedRealtime(t *testing.T) {
	t.Helper()
	h.wait(t, h.subscribedRealtime)
}

func (h *Hardware) wait(t *testing.T, ready chan struct{}) {
	t.Helper()
	select {
	case <-ready:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for hardware subscription")
	}
}
