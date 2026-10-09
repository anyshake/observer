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

func float32Word(value float32) [4]byte {
	var word [4]byte
	binary.LittleEndian.PutUint32(word[:], math.Float32bits(value))
	return word
}

func TestV3VariableDataFixtures(t *testing.T) {
	t.Parallel()

	protocol := &ExplorerProtoImplV3{
		ExplorerOptions: ExplorerOptions{Latitude: 11.5, Longitude: 22.5, Elevation: 33.5},
	}
	protocol.channelDataBuf = []ChannelData{{ChannelCode: "HNZ", Data: []int32{1}}}
	protocol.packetTimeObj = time.Unix(10, 0)
	protocol.collectedSamples = 9
	protocol.variableAllSet = true
	deviceID := uint32(1)
	protocol.deviceVariable.SetDeviceId(&deviceID)
	protocol.resetFlags()
	if len(protocol.channelDataBuf) != 0 || !protocol.packetTimeObj.IsZero() || protocol.collectedSamples != 0 {
		t.Fatalf("flags = %#v samples=%d", protocol.channelDataBuf, protocol.collectedSamples)
	}
	if !protocol.variableAllSet || protocol.GetDeviceId() == "N/A" {
		t.Fatal("resetFlags cleared device variables")
	}

	protocol.getVariableData(5000, 0, [4]byte{})
	if !protocol.variableAllSet {
		t.Fatal("packet with no variable bits was incomplete")
	}
	protocol.getVariableData(0, 0, [4]byte{0x78, 0x56, 0x34, 0x12})
	protocol.getVariableData(1000, 0, float32Word(1))
	protocol.getVariableData(2000, 0, float32Word(2))
	protocol.getVariableData(3000, 0, float32Word(3))
	protocol.getVariableData(4000, 0, float32Word(4))
	if protocol.GetDeviceId() != "N/A" {
		t.Fatalf("cleared device id = %s", protocol.GetDeviceId())
	}
	if _, err := protocol.GetTemperature(); err == nil {
		t.Fatal("cleared temperature was still readable")
	}
	lat, lon, elv, err := protocol.GetCoordinates(false)
	if err != nil || lat != 11.5 || lon != 22.5 || elv != 33.5 {
		t.Fatalf("fallback coordinates = %v %v %v, %v", lat, lon, elv, err)
	}

	protocol.resetVariables()
	if protocol.variableAllSet || protocol.GetDeviceId() != "N/A" {
		t.Fatal("resetVariables left device state behind")
	}
	const variableBits = uint32(0x1F)
	protocol.getVariableData(6000, variableBits, [4]byte{})
	if protocol.variableAllSet {
		t.Fatal("uncollected variable bits were treated as complete")
	}
	protocol.getVariableData(0, variableBits, [4]byte{0xAC, 0x81, 0x2F, 0x01})
	if protocol.variableAllSet || protocol.GetDeviceId() != "012F81AC" {
		t.Fatalf("device id = %s complete=%v", protocol.GetDeviceId(), protocol.variableAllSet)
	}
	protocol.getVariableData(1000, variableBits, float32Word(12.25))
	protocol.getVariableData(2000, variableBits, float32Word(23.5))
	protocol.getVariableData(3000, variableBits, float32Word(34.5))
	if protocol.variableAllSet {
		t.Fatal("variables settled before temperature")
	}
	if _, err = protocol.GetTemperature(); err == nil {
		t.Fatal("temperature was readable before its slot")
	}
	protocol.getVariableData(4000, variableBits, float32Word(18.5))
	if !protocol.variableAllSet {
		t.Fatal("variable bits did not settle")
	}
	lat, lon, elv, err = protocol.GetCoordinates(false)
	if err != nil || lat != float64(float32(12.25)) || lon != float64(float32(23.5)) || elv != float64(float32(34.5)) {
		t.Fatalf("GNSS coordinates = %v %v %v, %v", lat, lon, elv, err)
	}
	if got, err := protocol.GetTemperature(); err != nil || got != float64(float32(18.5)) {
		t.Fatalf("temperature = %v, %v", got, err)
	}

	protocol.deviceVariable.Reset()
	protocol.getVariableData(0, variableBits, [4]byte{0xFF, 0xFF, 0xFF, 0x7F})
	if protocol.GetDeviceId() != "FFFFFFFF" {
		t.Fatalf("saturated device id = %s", protocol.GetDeviceId())
	}
	if protocol.variableAllSet {
		t.Fatal("saturated device id completed the other variables")
	}
}

func TestV3OpenAndCloseErrors(t *testing.T) {
	t.Parallel()

	protocol := &ExplorerProtoImplV3{}
	if _, _, err := protocol.Open(context.Background()); err == nil || !strings.Contains(err.Error(), "transport") {
		t.Fatalf("nil transport error = %v", err)
	}
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
	if transport.opened != 1 || transport.flushed != 1 || transport.closed != 1 {
		t.Fatalf("flush cleanup opens=%d flushes=%d closes=%d", transport.opened, transport.flushed, transport.closed)
	}

	transport = &fakeTransport{}
	protocol.Transport = transport
	protocol.Logger = nil
	if _, _, err := protocol.Open(context.Background()); err == nil || !strings.Contains(err.Error(), "logger") {
		t.Fatalf("nil logger error = %v", err)
	}
	if transport.closed != 1 {
		t.Fatalf("logger cleanup closes = %d", transport.closed)
	}

	protocol.Logger = testLogger(nil)
	transport = &fakeTransport{}
	protocol.Transport = transport
	if _, _, err := protocol.Open(context.Background()); err == nil || !strings.Contains(err.Error(), "ntp") {
		t.Fatalf("empty NTP pool error = %v", err)
	}
	if transport.closed != 1 {
		t.Fatalf("NTP cleanup closes = %d", transport.closed)
	}

	if err := (&ExplorerProtoImplV3{}).Close(); err == nil {
		t.Fatal("nil transport close succeeded")
	}
	failing := &fakeTransport{closeErr: errors.New("close failed")}
	protocol = &ExplorerProtoImplV3{
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

func TestV3CommandSurface(t *testing.T) {
	t.Parallel()

	logs := &logBuffer{}
	protocol := &ExplorerProtoImplV3{
		Transport:          &fakeTransport{},
		Logger:             testLogger(logs),
		ExplorerOptions:    ExplorerOptions{Model: "missing-model"},
		messageBus:         message.NewBus[Event](EXPLORER_STREAM_TOPIC),
		messageBusRealtime: message.NewBus[Event](EXPLORER_REALTIME_STREAM_TOPIC),
	}
	if _, _, _, err := protocol.GetCoordinates(false); err == nil || !strings.Contains(err.Error(), "latitude") {
		t.Fatalf("missing latitude error = %v", err)
	}
	if _, err := protocol.GetMetadata("a", "b", "c", "d", "NET", "STA", "00", false); err == nil || !strings.Contains(err.Error(), "latitude") {
		t.Fatalf("metadata latitude error = %v", err)
	}
	latitude := 25.126
	protocol.deviceVariable.SetLatitude(&latitude)
	if _, _, _, err := protocol.GetCoordinates(false); err == nil || !strings.Contains(err.Error(), "longitude") {
		t.Fatalf("missing longitude error = %v", err)
	}
	if _, err := protocol.GetMetadata("a", "b", "c", "d", "NET", "STA", "00", false); err == nil || !strings.Contains(err.Error(), "longitude") {
		t.Fatalf("metadata longitude error = %v", err)
	}
	longitude := 121.234
	protocol.deviceVariable.SetLongitude(&longitude)
	if _, _, _, err := protocol.GetCoordinates(false); err == nil || !strings.Contains(err.Error(), "altitude") {
		t.Fatalf("missing altitude error = %v", err)
	}
	if _, err := protocol.GetTemperature(); err == nil || !strings.Contains(err.Error(), "temperature") {
		t.Fatalf("missing temperature error = %v", err)
	}
	if _, err := protocol.GetMetadata("a", "b", "c", "d", "NET", "STA", "00", false); err == nil || !strings.Contains(err.Error(), "altitude") {
		t.Fatalf("metadata error = %v", err)
	}
	elevation, temperature := 42.0, 18.5
	deviceID := uint32(0x012F81AC)
	protocol.deviceVariable.SetElevation(&elevation)
	if _, err := protocol.GetMetadata("a", "b", "c", "d", "NET", "STA", "00", true); err == nil || !strings.Contains(err.Error(), "missing-model") && !strings.Contains(err.Error(), "template") {
		t.Fatalf("missing model error = %v", err)
	}
	protocol.ExplorerOptions.Model = "E-C111G"
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
	protocol.deviceConfig.SetSampleRate(100)
	protocol.deviceConfig.SetPacketInterval(100 * time.Millisecond)
	protocol.deviceConfig.SetGnssAvailability(true)
	protocol.deviceConfig.SetChannelCodes([]string{"EHZ", "EHE"})
	protocol.deviceConfig.SetModel("E-C111G")
	protocol.deviceConfig.SetProtocol("v3")
	protocol.deviceStatus.SetStartedAt(started)
	protocol.deviceStatus.SetUpdatedAt(started)
	protocol.deviceStatus.IncrementFrames()
	protocol.deviceStatus.IncrementErrors()
	protocol.deviceStatus.IncrementMessages()
	if config := protocol.GetConfig(); config.GetSampleRate() != 100 || config.GetPacketInterval() != 100*time.Millisecond || config.GetProtocol() != "v3" {
		t.Fatalf("config snapshot = %#v", config)
	}
	if status := protocol.GetStatus(); !status.GetUpdatedAt().Equal(started) || status.GetFrames() != 1 || status.GetErrors() != 1 || status.GetMessages() != 1 {
		t.Fatalf("status snapshot frames=%d", status.GetFrames())
	}
	if render, err := protocol.GetMetadata("a", "b", "c", "d", "NET", "STA", "00", false); err != nil || render.SeisComP() == "" {
		t.Fatalf("metadata error = %v", err)
	}

	received := make(chan Event, 1)
	if err := protocol.Subscribe("client", nil); err == nil {
		t.Fatal("nil handler was accepted")
	}
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
	if err := protocol.SubscribeRealtime("boom-live", func(Event) { panic("boom") }); err != nil {
		t.Fatal(err)
	}
	protocol.messageBusRealtime.Publish(NewEvent(started, &protocol.deviceConfig, nil))
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) && (!strings.Contains(logs.String(), "normal stream subscriber") || !strings.Contains(logs.String(), "realtime stream subscriber")) {
		time.Sleep(5 * time.Millisecond)
	}
	if !strings.Contains(logs.String(), "normal stream subscriber") || !strings.Contains(logs.String(), "realtime stream subscriber") {
		t.Fatalf("panics were not reported: %s", logs.String())
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
	if err := protocol.SubscribeRealtime("after", func(Event) {}); err == nil {
		t.Fatal("realtime subscribe succeeded after close")
	}
}
