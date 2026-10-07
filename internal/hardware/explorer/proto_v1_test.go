package explorer

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/message"
	"github.com/rs/zerolog"
)

type fakeTransport struct {
	openErr, closeErr, flushErr, readErr error
	opened, closed, flushed              int
}

func (f *fakeTransport) Open() error {
	f.opened++
	return f.openErr
}
func (f *fakeTransport) Close() error {
	f.closed++
	return f.closeErr
}
func (f *fakeTransport) Flush() error {
	f.flushed++
	return f.flushErr
}
func (f *fakeTransport) Read(buf []byte) (int, error) {
	if f.readErr != nil {
		return 0, f.readErr
	}
	return 0, nil
}
func (f *fakeTransport) Write(buf []byte) (int, error)  { return len(buf), nil }
func (f *fakeTransport) GetLatency(int) time.Duration   { return 0 }
func (f *fakeTransport) SetTimeout(time.Duration) error { return nil }
func (f *fakeTransport) ReadUntil(context.Context, int, func(*[]byte, *time.Time) bool, time.Duration) ([]byte, bool, time.Duration, error) {
	if f.readErr != nil {
		return nil, false, 0, f.readErr
	}
	return nil, false, 0, nil
}

type logBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (l *logBuffer) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.buf.Write(p)
}

func (l *logBuffer) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.buf.String()
}

func testLogger(w *logBuffer) *logger.Adapter {
	if w == nil {
		return &logger.Adapter{Logger: zerolog.Nop()}
	}
	return &logger.Adapter{Logger: zerolog.New(w)}
}

func putInt32Samples(dst []byte, values []int32) {
	for i, value := range values {
		binary.LittleEndian.PutUint32(dst[i*4:], uint32(value))
	}
}

func buildV1DataPacket(headerSize, channelSize int, z, e, n []int32) []byte {
	zOffset := headerSize + int(binary.Size(int64(0)))
	eOffset := zOffset + channelSize*4
	nOffset := eOffset + channelSize*4
	packet := make([]byte, nOffset+channelSize*4+4)
	if headerSize > 0 {
		packet[0] = 0xFC
	}
	if headerSize > 1 {
		packet[1] = 0x1B
	}
	putInt32Samples(packet[zOffset:eOffset], z)
	putInt32Samples(packet[eOffset:nOffset], e)
	putInt32Samples(packet[nOffset:nOffset+channelSize*4], n)
	var checksum [3]byte
	for i := zOffset; i < eOffset; i++ {
		checksum[0] ^= packet[i]
	}
	for i := eOffset; i < nOffset; i++ {
		checksum[1] ^= packet[i]
	}
	for i := nOffset; i < len(packet)-4; i++ {
		checksum[2] ^= packet[i]
	}
	copy(packet[len(packet)-4:], checksum[:])
	return packet
}

func TestV1PacketLayout(t *testing.T) {
	t.Parallel()

	protocol := &ExplorerProtoImplV1{}
	if got := protocol.getPacketSize(2, 5); got != 66 {
		t.Fatalf("packet size = %d, want 66", got)
	}
	if got := protocol.getIndices(nil, []byte{0xFC, 0x1B}); len(got) != 0 {
		t.Fatalf("indices in empty buffer = %v", got)
	}
	if got := protocol.getIndices([]byte{0x00, 0xFC, 0x1B, 0xFC, 0x1B}, []byte{0xFC, 0x1B}); len(got) != 2 || got[0] != 1 || got[1] != 3 {
		t.Fatalf("header indices = %v, want [1 3]", got)
	}
	if got := protocol.getIndices([]byte{0xFC}, []byte{0xFC, 0x1B}); len(got) != 0 {
		t.Fatalf("short buffer indices = %v", got)
	}

	if _, err := protocol.fixSampleRate(5, 0); err == nil {
		t.Fatal("expected invalid duration")
	}
	rates := []struct {
		samples  int64
		duration time.Duration
		want     int
	}{
		{5, time.Millisecond, 5000},
		{5, 2 * time.Millisecond, 2500},
		{5, 4 * time.Millisecond, 1250},
		{1, time.Millisecond, 1000},
		{2, time.Millisecond, 2500},
		{6, time.Millisecond, 5000},
		{5, 8 * time.Millisecond, 500},
		{6, 8 * time.Millisecond, 500},
		{5, 10 * time.Millisecond, 500},
		{5, 20 * time.Millisecond, 250},
		{5, 40 * time.Millisecond, 125},
		{5, 100 * time.Millisecond, 50},
		{5, 200 * time.Millisecond, 50},
	}
	for _, rate := range rates {
		got, err := protocol.fixSampleRate(rate.samples, rate.duration)
		if err != nil || got != rate.want {
			t.Errorf("fixSampleRate(%d, %s) = %d, %v, want %d", rate.samples, rate.duration, got, err, rate.want)
		}
	}
}

func TestV1ChannelDataFixture(t *testing.T) {
	t.Parallel()

	protocol := &ExplorerProtoImplV1{ChannelCodes: []string{"HNZ"}}
	z := []int32{1, -2, 3, -4, 5}
	e := []int32{6, 7, 8, 9, 10}
	n := []int32{-11, 12, -13, 14, -15}
	packet := buildV1DataPacket(2, 5, z, e, n)
	got, err := protocol.getChannelData(packet, 2, 5)
	if err != nil {
		t.Fatal(err)
	}
	wantCodes := []string{"HNZ", "CH2", "CH3"}
	want := [][]int32{z, e, n}
	if len(got) != 3 {
		t.Fatalf("channels = %d, want 3", len(got))
	}
	for i, channel := range got {
		if channel.ChannelCode != wantCodes[i] || channel.ChannelId != i+1 || channel.ByteSize != 4 || channel.DataType != "int32" {
			t.Fatalf("channel %d metadata = %#v", i, channel)
		}
		if len(channel.Data) != len(want[i]) {
			t.Fatalf("channel %d samples = %v, want %v", i, channel.Data, want[i])
		}
		for j := range want[i] {
			if channel.Data[j] != want[i][j] {
				t.Fatalf("channel %d samples = %v, want %v", i, channel.Data, want[i])
			}
		}
	}
	if codes := protocol.deviceConfig.GetChannelCodes(); len(codes) != 3 || codes[0] != "HNZ" || codes[1] != "CH2" {
		t.Fatalf("stored channel codes = %#v", codes)
	}

	again, err := protocol.getChannelData(packet, 2, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(again[0].Data) != len(z) || len(protocol.channelDataBuf[0].Data) != len(z)*2 {
		t.Fatalf("repeat decode = %v, buffer %v", again[0].Data, protocol.channelDataBuf[0].Data)
	}

	packet[len(packet)-4] ^= 0xFF
	if _, err = protocol.getChannelData(packet, 2, 5); err == nil || !strings.Contains(err.Error(), "checksum mismatch") {
		t.Fatalf("corrupt checksum error = %v", err)
	}

	short := make([]byte, 22)
	if _, err = protocol.getChannelData(short, 2, 1); err == nil || !strings.Contains(err.Error(), "n-axis") {
		t.Fatalf("short packet error = %v", err)
	}

	emptyProtocol := &ExplorerProtoImplV1{}
	empty, err := emptyProtocol.getChannelData(make([]byte, 14), 2, 0)
	if err != nil || len(empty) != 3 || len(empty[2].Data) != 0 || empty[0].ChannelCode != "CH1" {
		t.Fatalf("empty channels = %#v, %v", empty, err)
	}
}

func TestV1OpenAndCloseErrors(t *testing.T) {
	t.Parallel()

	protocol := &ExplorerProtoImplV1{}
	if _, _, err := protocol.Open(context.Background()); err == nil || !strings.Contains(err.Error(), "transport") {
		t.Fatalf("nil transport error = %v", err)
	}
	transport := &fakeTransport{}
	protocol.Transport = transport
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

	protocol.Transport = transport
	if _, _, err := protocol.Open(context.Background()); err == nil || !strings.Contains(err.Error(), "ntp") {
		t.Fatalf("empty NTP pool error = %v", err)
	}
	if transport.opened != 1 || transport.closed != 1 {
		t.Fatalf("cleanup opens = %d closes = %d", transport.opened, transport.closed)
	}

	if err := (&ExplorerProtoImplV1{}).Close(); err == nil || !strings.Contains(err.Error(), "transport") {
		t.Fatal("nil transport close succeeded")
	}
	failing := &fakeTransport{closeErr: errors.New("close failed")}
	protocol = &ExplorerProtoImplV1{
		Transport:          failing,
		Logger:             testLogger(nil),
		messageBus:         message.NewBus[Event](EXPLORER_STREAM_TOPIC),
		messageBusRealtime: message.NewBus[Event](EXPLORER_REALTIME_STREAM_TOPIC),
	}
	if err := protocol.Close(); err == nil || !strings.Contains(err.Error(), "close failed") {
		t.Fatalf("close error = %v", err)
	}
	if err := protocol.Flush(); err != nil {
		t.Fatal(err)
	}
	failing.flushErr = errors.New("flush failed")
	if err := protocol.Flush(); err == nil || !strings.Contains(err.Error(), "flush failed") {
		t.Fatalf("flush error = %v", err)
	}
}

func TestV1CommandSurface(t *testing.T) {
	t.Parallel()

	logs := &logBuffer{}
	transport := &fakeTransport{}
	protocol := &ExplorerProtoImplV1{
		Transport:          transport,
		Logger:             testLogger(logs),
		ExplorerOptions:    ExplorerOptions{Model: "missing-model"},
		messageBus:         message.NewBus[Event](EXPLORER_STREAM_TOPIC),
		messageBusRealtime: message.NewBus[Event](EXPLORER_REALTIME_STREAM_TOPIC),
	}
	if _, _, _, err := protocol.GetCoordinates(false); err == nil || !strings.Contains(err.Error(), "latitude") {
		t.Fatalf("missing latitude error = %v", err)
	}
	if _, err := protocol.GetTemperature(); err == nil || !strings.Contains(err.Error(), "temperature") {
		t.Fatalf("missing temperature error = %v", err)
	}
	if got := protocol.GetDeviceId(); got != "N/A" {
		t.Fatalf("device id = %s", got)
	}
	if _, err := protocol.GetMetadata("a", "b", "c", "d", "NET", "STA", "00", false); err == nil {
		t.Fatal("metadata succeeded without coordinates")
	}

	latitude, longitude, elevation := 25.126, 121.234, 42.0
	temperature := 18.5
	deviceID := uint32(0x012F81AC)
	protocol.deviceVariable.SetLatitude(&latitude)
	if _, _, _, err := protocol.GetCoordinates(false); err == nil || !strings.Contains(err.Error(), "longitude") {
		t.Fatalf("missing longitude error = %v", err)
	}
	protocol.deviceVariable.SetLongitude(&longitude)
	if _, _, _, err := protocol.GetCoordinates(false); err == nil || !strings.Contains(err.Error(), "altitude") {
		t.Fatalf("missing altitude error = %v", err)
	}
	protocol.deviceVariable.SetElevation(&elevation)
	protocol.deviceVariable.SetTemperature(&temperature)
	protocol.deviceVariable.SetDeviceId(&deviceID)
	lat, lon, elv, err := protocol.GetCoordinates(true)
	if err != nil || lat != 25.13 || lon != 121.23 || elv != elevation {
		t.Fatalf("coordinates = %v %v %v, %v", lat, lon, elv, err)
	}
	if got, err := protocol.GetTemperature(); err != nil || got != temperature {
		t.Fatalf("temperature = %v, %v", got, err)
	}
	if got := protocol.GetDeviceId(); got != "012F81AC" {
		t.Fatalf("device id = %s", got)
	}

	started := time.Unix(1700000000, 0)
	protocol.deviceConfig.SetSampleRate(100)
	protocol.deviceConfig.SetPacketInterval(time.Second)
	protocol.deviceConfig.SetGnssAvailability(true)
	protocol.deviceConfig.SetChannelCodes([]string{"EHZ", "EHE", "EHN"})
	protocol.deviceConfig.SetModel("E-C111G")
	protocol.deviceConfig.SetProtocol("v1")
	protocol.deviceStatus.SetStartedAt(started)
	protocol.deviceStatus.SetUpdatedAt(started.Add(time.Second))
	protocol.deviceStatus.IncrementFrames()
	protocol.deviceStatus.IncrementErrors()
	protocol.deviceStatus.IncrementMessages()
	config := protocol.GetConfig()
	if config.GetSampleRate() != 100 || config.GetPacketInterval() != time.Second || !config.GetGnssAvailability() || config.GetProtocol() != "v1" || config.GetModel() != "E-C111G" {
		t.Fatalf("config snapshot = %#v", config)
	}
	status := protocol.GetStatus()
	if !status.GetStartedAt().Equal(started) || status.GetFrames() != 1 || status.GetErrors() != 1 || status.GetMessages() != 1 {
		t.Fatalf("status snapshot = %#v", status)
	}
	protocol.deviceStatus.IncrementFrames()
	if status.GetFrames() != 1 {
		t.Fatal("status snapshot aliased the live counters")
	}
	if render, err := protocol.GetMetadata("aff", "desc", "TW", "Taipei", "NET", "STA", "00", false); err == nil || render != nil {
		t.Fatal("missing metadata model succeeded")
	}
	protocol.ExplorerOptions.Model = "E-C111G"
	render, err := protocol.GetMetadata("aff", "desc", "TW", "Taipei", "NET", "STA", "00", true)
	if err != nil || render == nil || render.SeisComP() == "" || render.StationXML() == "" {
		t.Fatalf("metadata = %#v, %v", render, err)
	}

	received := make(chan Event, 1)
	handler := func(event Event) { received <- event }
	if err := protocol.Subscribe("client", nil); err == nil {
		t.Fatal("nil handler was accepted")
	}
	if err := protocol.Subscribe("client", handler); err != nil {
		t.Fatal(err)
	}
	if err := protocol.Subscribe("client", handler); err == nil {
		t.Fatal("duplicate subscription was accepted")
	}
	protocol.messageBus.Publish(NewEvent(started, &protocol.deviceConfig, nil))
	select {
	case event := <-received:
		if event.SampleRate != 100 || !event.GNSSAvailable {
			t.Fatalf("event = %#v", event)
		}
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
	protocol.messageBusRealtime.Publish(NewEvent(started, &protocol.deviceConfig, []ChannelData{{ChannelCode: "EHZ"}}))
	select {
	case event := <-realtime:
		if len(event.ChannelData) != 1 || event.ChannelData[0].ChannelCode != "EHZ" {
			t.Fatalf("realtime event = %#v", event)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for realtime event")
	}
	if err := protocol.SubscribeRealtime("panic", func(Event) { panic("boom") }); err != nil {
		t.Fatal(err)
	}
	protocol.messageBusRealtime.Publish(NewEvent(started, &protocol.deviceConfig, nil))
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) && !strings.Contains(logs.String(), "realtime stream subscriber") {
		time.Sleep(5 * time.Millisecond)
	}
	if !strings.Contains(logs.String(), "realtime stream subscriber") {
		t.Fatalf("panic was not reported: %s", logs.String())
	}
	if err := protocol.UnsubscribeRealtime("live"); err != nil {
		t.Fatal(err)
	}
	if err := protocol.UnsubscribeRealtime("missing"); err == nil {
		t.Fatal("missing realtime unsubscribe succeeded")
	}

	if err := protocol.Subscribe("panic-normal", func(Event) { panic("boom") }); err != nil {
		t.Fatal(err)
	}
	protocol.messageBus.Publish(NewEvent(started, &protocol.deviceConfig, nil))
	deadline = time.Now().Add(time.Second)
	for time.Now().Before(deadline) && !strings.Contains(logs.String(), "normal stream subscriber") {
		time.Sleep(5 * time.Millisecond)
	}
	if !strings.Contains(logs.String(), "normal stream subscriber") {
		t.Fatalf("archived panic was not reported: %s", logs.String())
	}

	if err := protocol.Close(); err != nil {
		t.Fatal(err)
	}
	if err := protocol.Subscribe("after-close", handler); err == nil {
		t.Fatal("subscribe succeeded after close")
	}
}
