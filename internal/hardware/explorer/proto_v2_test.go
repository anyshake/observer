package explorer

import (
	"context"
	"encoding/binary"
	"errors"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/anyshake/observer/pkg/message"
)

func buildV2DataPacket(channelSize int, z, e, n []int32) []byte {
	zOffset := 2 + 8 + 4
	eOffset := zOffset + channelSize*4
	nOffset := eOffset + channelSize*4
	packet := make([]byte, nOffset+channelSize*4+1)
	packet[0], packet[1] = 0xFA, 0xDE
	binary.LittleEndian.PutUint64(packet[2:10], 4_000)
	binary.LittleEndian.PutUint32(packet[10:14], 0x1234)
	putInt32Samples(packet[zOffset:eOffset], z)
	putInt32Samples(packet[eOffset:nOffset], e)
	putInt32Samples(packet[nOffset:nOffset+channelSize*4], n)
	return packet
}

func TestV2ChannelAndVariableFixtures(t *testing.T) {
	t.Parallel()

	protocol := &ExplorerProtoImplV2{
		ChannelCodes:    []string{"HNZ", "HNE"},
		ExplorerOptions: ExplorerOptions{Latitude: 11.5, Longitude: 22.5, Elevation: 33.5},
	}
	if got := protocol.getPacketSize(2, 5); got != 75 {
		t.Fatalf("packet size = %d, want 75", got)
	}

	z := []int32{1, -1}
	e := []int32{2, -2}
	n := []int32{3, -3}
	got, err := protocol.getChannelData(buildV2DataPacket(2, z, e, n), 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	wantCodes := []string{"HNZ", "HNE", "CH3"}
	want := [][]int32{z, e, n}
	for i, channel := range got {
		if channel.ChannelCode != wantCodes[i] || channel.ChannelId != i+1 || channel.ByteSize != 4 || channel.DataType != "int32" {
			t.Fatalf("channel %d = %#v", i, channel)
		}
		for j := range want[i] {
			if channel.Data[j] != want[i][j] {
				t.Fatalf("channel %d samples = %v, want %v", i, channel.Data, want[i])
			}
		}
	}
	if _, err = protocol.getChannelData(buildV2DataPacket(2, z, e, n), 2, 2); err != nil {
		t.Fatal(err)
	}
	if len(protocol.channelDataBuf[0].Data) != 4 {
		t.Fatalf("accumulated samples = %v", protocol.channelDataBuf[0].Data)
	}
	if codes := protocol.deviceConfig.GetChannelCodes(); len(codes) != 3 || codes[2] != "CH3" {
		t.Fatalf("channel codes = %#v", codes)
	}

	short := make([]byte, 23)
	if _, err = protocol.getChannelData(short, 2, 1); err == nil || !strings.Contains(err.Error(), "n-axis") {
		t.Fatalf("short packet error = %v", err)
	}
	empty, err := protocol.getChannelData(make([]byte, 15), 2, 0)
	if err != nil || len(empty) != 3 || len(empty[0].Data) != 0 {
		t.Fatalf("empty packet = %#v, %v", empty, err)
	}

	protocol.getVariableData(0, 0x1234)
	if !protocol.variableAllSet || protocol.deviceConfig.GetGnssAvailability() {
		t.Fatal("non-GNSS device word did not settle")
	}
	if gotID := protocol.GetDeviceId(); gotID != "00001234" {
		t.Fatalf("device id = %s", gotID)
	}
	if _, err = protocol.deviceVariable.GetLatitude(false); err == nil {
		t.Fatal("latitude was set before its variable slot")
	}
	protocol.getVariableData(1000, 0)
	protocol.getVariableData(2000, 0)
	if _, _, _, err = protocol.GetCoordinates(false); err == nil {
		t.Fatal("coordinates were complete before elevation")
	}
	protocol.getVariableData(3000, 0)
	lat, lon, elv, err := protocol.GetCoordinates(false)
	if err != nil || lat != 11.5 || lon != 22.5 || elv != 33.5 {
		t.Fatalf("fallback coordinates = %v %v %v, %v", lat, lon, elv, err)
	}

	protocol.resetVariables()
	if protocol.variableAllSet {
		t.Fatal("reset left variables complete")
	}
	if protocol.GetDeviceId() != "N/A" {
		t.Fatal("reset left the device id")
	}
	protocol.getVariableData(4000, 0x80000055)
	if protocol.variableAllSet || !protocol.deviceConfig.GetGnssAvailability() || protocol.GetDeviceId() != "00000055" {
		t.Fatalf("GNSS word = set:%v gnss:%v id:%s", protocol.variableAllSet, protocol.deviceConfig.GetGnssAvailability(), protocol.GetDeviceId())
	}
	protocol.getVariableData(5000, math.Float32bits(12.25))
	if _, err = protocol.deviceVariable.GetLongitude(false); err == nil {
		t.Fatal("longitude was accepted from the latitude slot")
	}
	protocol.getVariableData(6000, math.Float32bits(23.5))
	if _, err = protocol.deviceVariable.GetElevation(); err == nil {
		t.Fatal("elevation was accepted before its slot")
	}
	protocol.getVariableData(7000, math.Float32bits(34.5))
	if !protocol.variableAllSet {
		t.Fatal("GNSS variables did not settle")
	}
	lat, lon, elv, err = protocol.GetCoordinates(false)
	if err != nil || lat != float64(float32(12.25)) || lon != float64(float32(23.5)) || elv != float64(float32(34.5)) {
		t.Fatalf("GNSS coordinates = %v %v %v, %v", lat, lon, elv, err)
	}

	saturated := &ExplorerProtoImplV2{}
	saturated.getVariableData(8000, 0xFFFFFFFF)
	if gotID := saturated.GetDeviceId(); gotID != "FFFFFFFF" {
		t.Fatalf("saturated device id = %s", gotID)
	}
	if saturated.variableAllSet || !saturated.deviceConfig.GetGnssAvailability() {
		t.Fatal("saturated GNSS word did not keep variables incomplete")
	}
}

func TestV2OpenAndCloseErrors(t *testing.T) {
	t.Parallel()

	protocol := &ExplorerProtoImplV2{}
	if _, _, err := protocol.Open(context.Background()); err == nil || !strings.Contains(err.Error(), "transport") {
		t.Fatalf("nil transport error = %v", err)
	}
	protocol.Transport = &fakeTransport{}
	if _, _, err := protocol.Open(context.Background()); err == nil || !strings.Contains(err.Error(), "logger") {
		t.Fatalf("nil logger error = %v", err)
	}
	protocol.Logger = testLogger(nil)
	protocol.Transport = &fakeTransport{openErr: errors.New("open failed")}
	if _, _, err := protocol.Open(context.Background()); err == nil || !strings.Contains(err.Error(), "open failed") {
		t.Fatalf("open error = %v", err)
	}
	if protocol.Transport.(*fakeTransport).closed != 0 {
		t.Fatal("failed open closed the transport")
	}

	transport := &fakeTransport{flushErr: errors.New("flush failed")}
	protocol.Transport = transport
	if _, _, err := protocol.Open(context.Background()); err == nil || !strings.Contains(err.Error(), "flush failed") {
		t.Fatalf("flush error = %v", err)
	}
	if transport.opened != 1 || transport.closed != 1 || transport.flushed != 1 {
		t.Fatalf("flush cleanup opens=%d flushes=%d closes=%d", transport.opened, transport.flushed, transport.closed)
	}

	transport = &fakeTransport{}
	protocol.Transport = transport
	if _, _, err := protocol.Open(context.Background()); err == nil || !strings.Contains(err.Error(), "ntp") {
		t.Fatalf("empty NTP pool error = %v", err)
	}
	if transport.closed != 1 {
		t.Fatalf("NTP cleanup closes = %d", transport.closed)
	}

	if err := (&ExplorerProtoImplV2{}).Close(); err == nil {
		t.Fatal("nil transport close succeeded")
	}
	failing := &fakeTransport{closeErr: errors.New("close failed")}
	protocol = &ExplorerProtoImplV2{
		Transport:          failing,
		Logger:             testLogger(nil),
		messageBus:         message.NewBus[Event](EXPLORER_STREAM_TOPIC),
		messageBusRealtime: message.NewBus[Event](EXPLORER_REALTIME_STREAM_TOPIC),
	}
	if err := protocol.Close(); err == nil || !strings.Contains(err.Error(), "close failed") {
		t.Fatalf("close error = %v", err)
	}
	failing.flushErr = errors.New("flush failed")
	if err := protocol.Flush(); err == nil || !strings.Contains(err.Error(), "flush failed") {
		t.Fatalf("flush error = %v", err)
	}
}

func TestV2CommandSurface(t *testing.T) {
	t.Parallel()

	logs := &logBuffer{}
	protocol := &ExplorerProtoImplV2{
		Transport:          &fakeTransport{},
		Logger:             testLogger(logs),
		ExplorerOptions:    ExplorerOptions{Model: "missing-model"},
		messageBus:         message.NewBus[Event](EXPLORER_STREAM_TOPIC),
		messageBusRealtime: message.NewBus[Event](EXPLORER_REALTIME_STREAM_TOPIC),
	}
	if _, _, _, err := protocol.GetCoordinates(false); err == nil || !strings.Contains(err.Error(), "latitude") {
		t.Fatalf("missing latitude error = %v", err)
	}
	latitude := 25.126
	protocol.deviceVariable.SetLatitude(&latitude)
	if _, _, _, err := protocol.GetCoordinates(false); err == nil || !strings.Contains(err.Error(), "longitude") {
		t.Fatalf("missing longitude error = %v", err)
	}
	longitude := 121.234
	protocol.deviceVariable.SetLongitude(&longitude)
	if _, _, _, err := protocol.GetCoordinates(false); err == nil || !strings.Contains(err.Error(), "altitude") {
		t.Fatalf("missing altitude error = %v", err)
	}
	if _, err := protocol.GetTemperature(); err == nil || !strings.Contains(err.Error(), "temperature") {
		t.Fatalf("missing temperature error = %v", err)
	}
	if protocol.GetDeviceId() != "N/A" {
		t.Fatal("expected placeholder device id")
	}
	if _, err := protocol.GetMetadata("a", "b", "c", "d", "NET", "STA", "00", false); err == nil {
		t.Fatal("metadata succeeded without coordinates")
	}

	elevation, temperature := 42.0, 18.5
	deviceID := uint32(0x012F81AC)
	protocol.deviceVariable.SetElevation(&elevation)
	protocol.deviceVariable.SetTemperature(&temperature)
	protocol.deviceVariable.SetDeviceId(&deviceID)
	if lat, lon, elv, err := protocol.GetCoordinates(true); err != nil || lat != 25.13 || lon != 121.23 || elv != elevation {
		t.Fatalf("coordinates = %v %v %v, %v", lat, lon, elv, err)
	}
	if got, err := protocol.GetTemperature(); err != nil || got != temperature {
		t.Fatalf("temperature = %v, %v", got, err)
	}
	if protocol.GetDeviceId() != "012F81AC" {
		t.Fatalf("device id = %s", protocol.GetDeviceId())
	}
	started := time.Unix(1700000000, 0)
	protocol.deviceConfig.SetSampleRate(50)
	protocol.deviceConfig.SetPacketInterval(200 * time.Millisecond)
	protocol.deviceConfig.SetGnssAvailability(true)
	protocol.deviceConfig.SetChannelCodes([]string{"HNZ"})
	protocol.deviceConfig.SetModel("E-C111G")
	protocol.deviceConfig.SetProtocol("v2")
	protocol.deviceStatus.SetStartedAt(started)
	protocol.deviceStatus.IncrementFrames()
	protocol.deviceStatus.IncrementErrors()
	protocol.deviceStatus.IncrementMessages()
	if config := protocol.GetConfig(); config.GetSampleRate() != 50 || !config.GetGnssAvailability() || config.GetProtocol() != "v2" {
		t.Fatalf("config snapshot = %#v", config)
	}
	if status := protocol.GetStatus(); status.GetFrames() != 1 || status.GetErrors() != 1 || status.GetMessages() != 1 || !status.GetStartedAt().Equal(started) {
		t.Fatalf("status snapshot frames=%d", status.GetFrames())
	}
	if render, err := protocol.GetMetadata("a", "b", "c", "d", "NET", "STA", "00", false); err == nil || render != nil {
		t.Fatal("missing model succeeded")
	}
	protocol.ExplorerOptions.Model = "E-C111G"
	if render, err := protocol.GetMetadata("a", "b", "c", "d", "NET", "STA", "00", false); err != nil || render.StationXML() == "" {
		t.Fatalf("metadata error = %v", err)
	}

	if err := protocol.Subscribe("client", nil); err == nil {
		t.Fatal("nil handler was accepted")
	}
	received := make(chan Event, 1)
	if err := protocol.Subscribe("client", func(event Event) { received <- event }); err != nil {
		t.Fatal(err)
	}
	if err := protocol.Subscribe("client", func(Event) {}); err == nil {
		t.Fatal("duplicate subscription was accepted")
	}
	protocol.messageBus.Publish(NewEvent(started, &protocol.deviceConfig, nil))
	select {
	case <-received:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for archived event")
	}
	if err := protocol.Unsubscribe("client"); err != nil {
		t.Fatal(err)
	}
	if err := protocol.Unsubscribe("missing"); err == nil {
		t.Fatal("missing unsubscribe succeeded")
	}
	realtime := make(chan Event, 1)
	if err := protocol.SubscribeRealtime("live", func(event Event) { realtime <- event }); err != nil {
		t.Fatal(err)
	}
	protocol.messageBusRealtime.Publish(NewEvent(started, &protocol.deviceConfig, nil))
	select {
	case <-realtime:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for realtime event")
	}
	if err := protocol.Subscribe("boom", func(Event) { panic("boom") }); err != nil {
		t.Fatal(err)
	}
	protocol.messageBus.Publish(NewEvent(started, &protocol.deviceConfig, nil))
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) && !strings.Contains(logs.String(), "normal stream subscriber") {
		time.Sleep(5 * time.Millisecond)
	}
	if !strings.Contains(logs.String(), "normal stream subscriber") {
		t.Fatalf("panic was not reported: %s", logs.String())
	}
	if err := protocol.SubscribeRealtime("boom-live", func(Event) { panic("boom") }); err != nil {
		t.Fatal(err)
	}
	protocol.messageBusRealtime.Publish(NewEvent(started, &protocol.deviceConfig, nil))
	deadline = time.Now().Add(time.Second)
	for time.Now().Before(deadline) && !strings.Contains(logs.String(), "realtime stream subscriber") {
		time.Sleep(5 * time.Millisecond)
	}
	if !strings.Contains(logs.String(), "realtime stream subscriber") {
		t.Fatalf("realtime panic was not reported: %s", logs.String())
	}
	if err := protocol.UnsubscribeRealtime("live"); err != nil {
		t.Fatal(err)
	}
	if err := protocol.UnsubscribeRealtime("missing"); err == nil {
		t.Fatal("missing realtime unsubscribe succeeded")
	}
	if err := protocol.Close(); err != nil {
		t.Fatal(err)
	}
	if err := protocol.Subscribe("after", func(Event) {}); err == nil {
		t.Fatal("subscribe succeeded after close")
	}
}
