package winston

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/ringbuf"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/bclswl0827/winsgo"
)

func TestLifecycle(t *testing.T) {
	logger.Init()
	_, handler := testsupport.OpenDAO(t)
	hw := testsupport.NewHardware()
	fixed := time.Date(2026, time.October, 8, 1, 2, 3, 0, time.UTC)
	svc := New(hw, handler, timesource.New(func() time.Time { return fixed }))
	if err := svc.Start(); err == nil {
		t.Fatal("start succeeded before initialization")
	}
	if svc.GetName() != "Winston Service" || svc.GetDescription() == "" {
		t.Fatal("service metadata is incomplete")
	}
	if _, err := svc.GetAssetList(); err == nil {
		t.Fatal("asset list succeeded")
	}
	if _, err := svc.GetAssetData("missing"); err == nil {
		t.Fatal("asset data succeeded")
	}
	if err := svc.Init(); err == nil {
		t.Fatal("Init succeeded without station settings")
	}

	testsupport.InitStation(t, handler)
	if err := svc.Init(); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "listen_host", "127.0.0.1"); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "listen_port", int64(freePort(t))); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "enabled", true); err != nil {
		t.Fatal(err)
	}
	if err := svc.Init(); err != nil {
		t.Fatal(err)
	}
	if !svc.IsEnabled() || svc.GetListenPort() == 0 {
		t.Fatal("service was not configured")
	}

	hw.SubscribeErr = errors.New("subscribe failed")
	if err := svc.Start(); err == nil {
		t.Fatal("subscribe failure was ignored")
	}
	hw.SubscribeErr = nil
	if err := svc.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.Stop() })
	hw.WaitSubscribed(t)
	hw.Emit(explorer.Event{
		Timestamp:  fixed,
		SampleRate: 100,
		ChannelData: []explorer.ChannelData{{
			ChannelCode: "EHZ",
			Data:        []int32{1, 2, 3, 4},
		}},
	})
	if err := svc.Start(); err == nil {
		t.Fatal("second start was accepted")
	}
	if err := svc.Restart(); err != nil {
		t.Fatal(err)
	}
	hw.WaitSubscribed(t)
	if svc.GetStatus().GetRestarts() != 1 {
		t.Fatalf("restarts = %d", svc.GetStatus().GetRestarts())
	}
	if err := svc.Stop(); err != nil {
		t.Fatal(err)
	}
}

func TestWaveformQueries(t *testing.T) {
	logger.Init()
	database, handler := testsupport.OpenDAO(t)
	testsupport.InitStation(t, handler)
	hw := testsupport.NewHardware()
	fixed := time.Date(2026, time.October, 8, 1, 2, 3, 0, time.UTC)
	svc := New(hw, handler, timesource.New(func() time.Time { return fixed }))
	if err := svc.Init(); err != nil {
		t.Fatal(err)
	}
	provider := &provider{
		hardwareDev: hw, timeSource: func() time.Time { return fixed },
		stationCode: svc.stationCode, networkCode: svc.networkCode, locationCode: svc.locationCode,
	}
	channels, err := provider.Channels(context.Background())
	if err != nil || len(channels) != 3 {
		t.Fatalf("channels = %#v, %v", channels, err)
	}
	if _, err := provider.Channels(canceledContext()); err == nil {
		t.Fatal("canceled provider succeeded")
	}

	buffer := ringbuf.New[winstonRingBuffer](8)
	samples := make([]int32, 100)
	buffer.Push(winstonRingBuffer{
		timestamp: fixed, sampleRate: 100,
		channelData: cloneChannelData([]explorer.ChannelData{{ChannelCode: "EHZ", Data: samples}}),
	})
	reader := &consumer{
		stationCode: svc.stationCode, networkCode: svc.networkCode, locationCode: svc.locationCode,
		ringBuffer: buffer, store: handler,
	}
	request := winsgo.WaveformRequest{
		Channel:   winsgo.SCNL{Station: svc.stationCode, Network: svc.networkCode, Location: svc.locationCode, Channel: "EHZ"},
		StartTime: fixed, EndTime: fixed.Add(500 * time.Millisecond),
	}
	if err := reader.Consume(context.Background(), request, func(waveform winsgo.Waveform) error {
		if len(waveform.Samples) == 0 {
			t.Fatal("empty waveform")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	request.Channel.Station = "OTHER"
	if !errors.Is(reader.Consume(context.Background(), request, nil), winsgo.ErrNoData) {
		t.Fatal("wrong station returned data")
	}
	request.Channel.Station = svc.stationCode
	request.EndTime = request.StartTime
	if !errors.Is(reader.Consume(context.Background(), request, nil), winsgo.ErrNoData) {
		t.Fatal("empty range returned data")
	}
	request.EndTime = request.StartTime.Add(2 * time.Hour)
	if err := reader.Consume(context.Background(), request, nil); err == nil {
		t.Fatal("oversized query was accepted")
	}

	empty := &consumer{stationCode: svc.stationCode, networkCode: svc.networkCode, locationCode: svc.locationCode, store: handler}
	request.EndTime = request.StartTime.Add(time.Second)
	if err := empty.Consume(canceledContext(), request, nil); err == nil {
		t.Fatal("canceled query succeeded")
	}
	if _, _, err := buildWaveform(fixed, fixed.Add(time.Second), []waveformRecord{{sampleRate: 0, samples: samples}}); err == nil {
		t.Fatal("invalid sample rate was accepted")
	}
	if _, _, err := buildWaveform(fixed, fixed.Add(time.Second), []waveformRecord{
		{timestamp: fixed, sampleRate: 100, samples: samples},
		{timestamp: fixed, sampleRate: 50, samples: samples},
	}); err == nil {
		t.Fatal("mixed sample rate was accepted")
	}

	testsupport.EnsureSeisTable(t, database, fixed)
	record := model.SeisRecord{}
	if err := record.Encode(fixed, 100, []explorer.ChannelData{{ChannelCode: "EHZ", Data: samples}}); err != nil {
		t.Fatal(err)
	}
	if err := handler.SeisRecordsCreate(record); err != nil {
		t.Fatal(err)
	}
	if err := empty.Consume(context.Background(), request, func(winsgo.Waveform) error { return nil }); err != nil {
		t.Fatal(err)
	}
}

func setConstraint(t *testing.T, handler *action.Handler, svc *WinstonServiceImpl, key string, value any) error {
	t.Helper()
	for _, constraint := range svc.GetConfigConstraint() {
		if constraint.GetKey() == key {
			return constraint.Set(handler, value)
		}
	}
	t.Fatalf("constraint %s was not found", key)
	return nil
}

func freePort(t *testing.T) int {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

func canceledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}
