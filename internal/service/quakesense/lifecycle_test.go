package quakesense

import (
	"context"
	"errors"
	"math"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/notification"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/timesource"
)

func TestFiltersAndLifecycle(t *testing.T) {
	logger.Init()
	_, handler := testsupport.OpenDAO(t)
	testsupport.InitStation(t, handler)
	hw := testsupport.NewHardware()
	fixed := time.Date(2026, time.October, 8, 1, 2, 3, 0, time.UTC)
	hub := notification.NewHub(4)
	events, unsubscribe := hub.Subscribe()
	defer unsubscribe()
	svc := New(t.Context(), hw, handler, timesource.New(func() time.Time { return fixed }), hub)
	if svc.GetName() == "" || svc.GetDescription() == "" || svc.GetStatus() == nil {
		t.Fatal("service metadata is incomplete")
	}
	if svc.IsEnabled(t.Context()) {
		t.Fatal("enabled before initialization")
	}
	if _, err := svc.GetAssetList(); err == nil {
		t.Fatal("asset list succeeded")
	}
	if _, err := svc.GetAssetData("missing"); err == nil {
		t.Fatal("asset data succeeded")
	}
	if err := svc.validateFinite("sample", math.NaN()); err == nil {
		t.Fatal("NaN was accepted")
	}
	if _, err := svc.newFilter(0); err == nil {
		t.Fatal("zero sample rate was accepted")
	}
	svc.filterType = "unknown"
	if _, err := svc.newFilter(100); err == nil {
		t.Fatal("unknown filter was accepted")
	}
	svc.filterType = NO_FILTER
	if filter, err := svc.newFilter(100); err != nil || filter != nil {
		t.Fatalf("no filter = %#v, %v", filter, err)
	}
	svc.minFreq, svc.maxFreq = 0.5, 10
	for _, kind := range []string{BAND_PASS_FILTER, LOW_PASS_FILTER, HIGH_PASS_FILTER} {
		svc.filterType = kind
		if _, err := svc.newFilter(100); err != nil {
			t.Fatalf("%s filter: %v", kind, err)
		}
	}
	svc.triggerMethod = CLASSIC_STA_LTA
	svc.ltaWindow, svc.staWindow = 0.2, 0.05
	if svc.filterBufferSize(100) < 1 {
		t.Fatal("empty classic buffer")
	}
	svc.triggerMethod = DELAYED_STA_LTA
	if svc.filterBufferSize(100) < 50 {
		t.Fatal("delayed buffer is too small")
	}

	broker := serveMQTT(t)
	if err := setConstraint(t, handler, svc, "mqtt_broker", broker); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "mqtt_username", "observer"); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "mqtt_password", "secret"); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "filter_type", NO_FILTER); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "sta_window", 0.05); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "lta_window", 0.2); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "trig_off", 0.2); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "trig_on", 0.5); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "enabled", true); err != nil {
		t.Fatal(err)
	}
	if err := svc.Init(t.Context()); err != nil {
		t.Fatal(err)
	}
	if !svc.IsEnabled(t.Context()) {
		t.Fatal("service stayed disabled")
	}
	if err := svc.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { testsupport.StopService(t, svc) })
	hw.WaitSubscribedRealtime(t)
	if err := svc.Start(t.Context()); err == nil {
		t.Fatal("duplicate Start succeeded")
	}
	if err := svc.Init(t.Context()); err == nil {
		t.Fatal("Init succeeded while the service was running")
	}

	stamp := fixed
	emit := func(rate int, gnss bool, samples []int32) {
		hw.EmitRealtime(explorer.Event{
			Timestamp: stamp, SampleRate: rate, GNSSAvailable: gnss,
			ChannelData: []explorer.ChannelData{{ChannelCode: "EHZ", Data: samples}},
		})
		if rate > 0 && len(samples) > 0 {
			stamp = stamp.Add(time.Duration(len(samples)) * time.Second / time.Duration(rate))
		}
	}
	hw.EmitRealtime(explorer.Event{Timestamp: stamp, SampleRate: 100, ChannelData: []explorer.ChannelData{{ChannelCode: "EHN", Data: []int32{1}}}})
	emit(100, false, nil)
	emit(0, false, []int32{1})
	emit(100, true, samples(40, 0))
	emit(100, false, samples(40, 100000))
	stamp = stamp.Add(time.Minute)
	emit(100, false, samples(40, 0))

	for _, method := range []string{DELAYED_STA_LTA, RECURSIVE_STA_LTA, Z_DETECT, "unknown", CLASSIC_STA_LTA} {
		svc.triggerMethod = method
		emit(90, false, samples(80, 100000))
	}
	svc.filterType = LOW_PASS_FILTER
	emit(80, false, samples(80, 1))
	svc.filterType = HIGH_PASS_FILTER
	emit(70, false, samples(80, 1))
	svc.filterType = BAND_PASS_FILTER
	emit(60, false, samples(80, 1))
	svc.filterType = NO_FILTER
	emit(50, false, samples(80, 1))

	select {
	case <-events:
	case <-time.After(200 * time.Millisecond):
	}
	if err := svc.Restart(t.Context()); err != nil {
		t.Fatal(err)
	}
	hw.WaitSubscribedRealtime(t)
	if svc.GetStatus().GetRestarts() != 1 {
		t.Fatalf("restarts = %d", svc.GetStatus().GetRestarts())
	}
	if err := svc.Stop(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func TestInitRejectsInvertedWindows(t *testing.T) {
	_, handler := testsupport.OpenDAO(t)
	testsupport.InitStation(t, handler)
	svc := New(t.Context(), testsupport.NewHardware(), handler, timesource.New(nil), nil)
	if err := svc.Init(t.Context()); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "sta_window", 20.0); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "lta_window", 4.0); err != nil {
		t.Fatal(err)
	}
	if err := svc.Init(t.Context()); err == nil {
		t.Fatal("LTA window shorter than STA was accepted")
	}
}

func TestStartConnectHonorsContext(t *testing.T) {
	svc := New(t.Context(), testsupport.NewHardware(), nil, timesource.New(nil), nil)
	svc.mqttBroker = serveSilentMQTT(t)
	svc.mqttClientId = "observer-test"

	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()
	started := time.Now()
	if err := svc.Start(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Start with an unresponsive broker = %v", err)
	}
	if elapsed := time.Since(started); elapsed >= quakeSenseConnectTimeout {
		t.Fatalf("Start returned after %s, ignoring the caller context", elapsed)
	}
	if svc.GetStatus().GetIsRunning() {
		t.Fatal("service is running after a failed MQTT connection")
	}
}

func serveSilentMQTT(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	var (
		mu    sync.Mutex
		conns []net.Conn
	)
	t.Cleanup(func() {
		_ = listener.Close()
		mu.Lock()
		defer mu.Unlock()
		for _, conn := range conns {
			_ = conn.Close()
		}
	})
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			mu.Lock()
			conns = append(conns, conn)
			mu.Unlock()
		}
	}()
	return "tcp://" + listener.Addr().String()
}

type blockingCoordinatesHardware struct {
	*testsupport.Hardware
	entered chan struct{}
	release chan struct{}
}

func (h *blockingCoordinatesHardware) GetCoordinates(fuzzy bool) (float64, float64, float64, error) {
	select {
	case h.entered <- struct{}{}:
	default:
	}
	<-h.release
	return h.Hardware.GetCoordinates(fuzzy)
}

func TestStopDoesNotWaitForEventPublishing(t *testing.T) {
	logger.Init()
	hw := &blockingCoordinatesHardware{
		Hardware: testsupport.NewHardware(),
		entered:  make(chan struct{}, 1),
		release:  make(chan struct{}),
	}
	svc := New(t.Context(), hw, nil, timesource.New(nil), nil)
	svc.mqttBroker = serveMQTT(t)
	svc.mqttClientId = "observer-test"
	svc.mqttTopic = "observer/test"
	svc.monitorChannel = "EHZ"
	svc.triggerMethod = CLASSIC_STA_LTA
	svc.filterType = NO_FILTER
	svc.staWindow, svc.ltaWindow = 0.05, 0.2
	svc.trigOn, svc.trigOff = 0.5, 0.2
	if err := svc.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { testsupport.StopService(t, svc) })
	hw.WaitSubscribedRealtime(t)

	stamp := time.Now()
	emitted := make(chan struct{})
	go func() {
		defer close(emitted)
		hw.EmitRealtime(explorer.Event{Timestamp: stamp, SampleRate: 100, ChannelData: []explorer.ChannelData{{ChannelCode: "EHZ", Data: samples(40, 0)}}})
		hw.EmitRealtime(explorer.Event{Timestamp: stamp.Add(400 * time.Millisecond), SampleRate: 100, ChannelData: []explorer.ChannelData{{ChannelCode: "EHZ", Data: samples(40, 100000)}}})
	}()
	select {
	case <-hw.entered:
	case <-time.After(2 * time.Second):
		close(hw.release)
		t.Fatal("no seismic event was detected")
	}

	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	err := svc.Stop(ctx)
	close(hw.release)
	<-emitted
	if err != nil {
		t.Fatalf("Stop waited for event publishing: %v", err)
	}
}

func samples(count int, value int32) []int32 {
	data := make([]int32, count)
	for i := range data {
		data[i] = value
	}
	return data
}

func serveMQTT(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = listener.Close() })
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go func() {
				defer conn.Close()
				buf := make([]byte, 4096)
				if _, err := conn.Read(buf); err != nil {
					return
				}
				_, _ = conn.Write([]byte{0x20, 0x02, 0x00, 0x00})
				for {
					n, err := conn.Read(buf)
					if err != nil {
						return
					}
					if n > 0 && buf[0] == 0xC0 {
						_, _ = conn.Write([]byte{0xD0, 0x00})
					}
				}
			}()
		}
	}()
	return "tcp://" + listener.Addr().String()
}

func setConstraint(t *testing.T, handler *action.Handler, svc *QuakeSenseServiceImpl, key string, value any) error {
	t.Helper()
	for _, constraint := range svc.GetConfigConstraint() {
		if constraint.GetKey() == key {
			return constraint.Set(t.Context(), handler, value)
		}
	}
	t.Fatalf("constraint %s was not found", key)
	return nil
}
