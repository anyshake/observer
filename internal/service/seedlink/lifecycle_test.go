package seedlink

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/message"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/bclswl0827/slgo/handlers"
)

func TestProviderAndConsumer(t *testing.T) {
	logger.Init()
	database, handler := testsupport.OpenDAO(t)
	testsupport.InitStation(t, handler)
	hw := testsupport.NewHardware()
	fixed := time.Date(2026, time.October, 8, 1, 2, 3, 0, time.UTC)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	source := &provider{
		ctx: ctx, hardwareDev: hw, timeSource: timesource.New(func() time.Time { return fixed }),
		actionHandler: handler, startTime: fixed, stationCode: "STA", networkCode: "NET", locationCode: "00",
	}
	if source.GetSoftware() == "" || source.GetOrganization() == "" || len(source.GetCapabilities()) == 0 {
		t.Fatal("provider metadata is incomplete")
	}
	if !source.GetCurrentTime().Equal(fixed) || !source.GetStartTime().Equal(fixed) {
		t.Fatal("provider clock mismatch")
	}
	stations := source.GetStations()
	if len(stations) != 1 || stations[0].Station != "STA" || stations[0].Network != "NET" {
		t.Fatalf("stations = %#v", stations)
	}
	streams := source.GetStreams()
	if len(streams) != 3 || streams[0].SeedName == "" {
		t.Fatalf("streams = %#v", streams)
	}
	if packets, err := source.QueryHistory(fixed, time.Time{}, nil); err != nil || packets != nil {
		t.Fatalf("empty query = %#v, %v", packets, err)
	}

	testsupport.EnsureSeisTable(t, database, fixed)
	record := model.SeisRecord{}
	if err := record.Encode(fixed, 100, []explorer.ChannelData{{ChannelCode: "EHZ", Data: []int32{1, 2, 3}}}); err != nil {
		t.Fatal(err)
	}
	if err := handler.SeisRecordsCreate(record); err != nil {
		t.Fatal(err)
	}
	packets, err := source.QueryHistory(fixed.Add(-time.Second), fixed.Add(time.Second), []handlers.SeedLinkChannel{{ChannelName: "EHZ"}, {ChannelName: "EHN"}})
	if err != nil || len(packets) != 1 || packets[0].Channel != "EHZ" {
		t.Fatalf("history = %#v, %v", packets, err)
	}
	canceled, stop := context.WithCancel(context.Background())
	stop()
	source.ctx = canceled
	if _, err := source.QueryHistory(fixed.Add(-time.Second), fixed.Add(time.Second), []handlers.SeedLinkChannel{{ChannelName: "EHZ"}}); err == nil {
		t.Fatal("canceled history query succeeded")
	}

	left, right := net.Pipe()
	defer left.Close()
	defer right.Close()
	client := &handlers.SeedLinkClient{Conn: left}
	registry := newSeedLinkClientRegistry()
	registry.Add(client)
	if _, ok := registry.Get(client.RemoteAddr().String()); !ok {
		t.Fatal("client was not registered")
	}
	registry.Remove(client)
	if _, ok := registry.Get(client.RemoteAddr().String()); ok {
		t.Fatal("client stayed registered")
	}
	registry.Add(client)
	hooks := &hooks{clients: registry}
	hooks.OnConnection(client)
	hooks.OnData(client, []byte("DATA"))
	hooks.OnCommand(client, []string{"HELLO"})
	hooks.OnClose(client)

	bus := message.NewBus[explorer.Event](ID)
	defer bus.Close()
	registry.Add(client)
	consumer := &consumer{messageBus: bus, clients: registry}
	delivered := make(chan handlers.SeedLinkDataPacket, 1)
	if err := consumer.Subscribe(client.RemoteAddr().String(), []handlers.SeedLinkChannel{{ChannelName: "EHZ"}}, func(packet handlers.SeedLinkDataPacket) {
		delivered <- packet
	}); err != nil {
		t.Fatal(err)
	}
	bus.Publish(explorer.Event{
		Timestamp: fixed, SampleRate: 100,
		ChannelData: []explorer.ChannelData{{ChannelCode: "EHN", Data: []int32{1}}, {ChannelCode: "EHZ", Data: []int32{4, 5}}},
	})
	select {
	case packet := <-delivered:
		if packet.Channel != "EHZ" || len(packet.DataArr) != 2 {
			t.Fatalf("packet = %#v", packet)
		}
	case <-time.After(time.Second):
		t.Fatal("realtime sample was not delivered")
	}
	registry.Remove(client)
	bus.Publish(explorer.Event{Timestamp: fixed, SampleRate: 100, ChannelData: []explorer.ChannelData{{ChannelCode: "EHZ", Data: []int32{1}}}})
	if err := consumer.Unsubscribe(client.RemoteAddr().String()); err != nil {
		t.Fatal(err)
	}
}

func TestLifecycle(t *testing.T) {
	logger.Init()
	_, handler := testsupport.OpenDAO(t)
	testsupport.InitStation(t, handler)
	hw := testsupport.NewHardware()
	fixed := time.Date(2026, time.October, 8, 1, 2, 3, 0, time.UTC)
	svc := New(hw, handler, timesource.New(func() time.Time { return fixed }))
	if svc.GetName() != "SeedLink Service" || svc.GetDescription() == "" || svc.GetStatus() == nil {
		t.Fatal("service metadata is incomplete")
	}
	if _, err := svc.GetAssetList(); err == nil {
		t.Fatal("asset list succeeded")
	}
	if _, err := svc.GetAssetData("missing"); err == nil {
		t.Fatal("asset data succeeded")
	}
	if err := svc.Init(); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "listen_host", "127.0.0.1"); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "listen_port", int64(freePort(t))); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "use_compression", true); err != nil {
		t.Fatal(err)
	}
	if err := setConstraint(t, handler, svc, "compress_scheme", COMPRESS_SCHEME_STEIM_2); err != nil {
		t.Fatal(err)
	}
	if err := svc.Init(); err != nil {
		t.Fatal(err)
	}
	if svc.GetListenPort() == 0 || !svc.IsEnabled() {
		t.Fatal("service was not configured")
	}
	if err := svc.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.Stop() })
	hw.WaitSubscribed(t)
	hw.Emit(explorer.Event{Timestamp: fixed, SampleRate: 100, ChannelData: []explorer.ChannelData{{ChannelCode: "EHZ", Data: []int32{1}}}})
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(svc.listenHost, strconv.Itoa(svc.listenPort)), time.Second)
	if err == nil {
		_ = conn.Close()
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

func setConstraint(t *testing.T, handler *action.Handler, svc *SeedLinkServiceImpl, key string, value any) error {
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
