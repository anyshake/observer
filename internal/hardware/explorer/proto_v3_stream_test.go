package explorer

import (
	"context"
	"errors"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/anyshake/observer/pkg/fifo"
	"github.com/anyshake/observer/pkg/message"
	"github.com/anyshake/observer/pkg/timesource"
)

func TestV3OpenGNSSStream(t *testing.T) {
	transport := newScriptTransport()
	base := time.Unix(1_700_000_000, 0)
	protocol := &ExplorerProtoImplV3{
		ChannelCodes:    []string{"HNZ", "HNE", "HNN"},
		ExplorerOptions: ExplorerOptions{Protocol: "v3", Model: "E-C111G", Latitude: 25, Longitude: 121, Elevation: 10},
		NtpOptions:      NtpOptions{Pool: []string{startTestNTP(t)}, Retry: 0, ReadTimeout: 1},
		Logger:          testLogger(nil),
		TimeSource:      timesource.New(func() time.Time { return base }),
		Transport:       transport,
	}
	t.Cleanup(func() { _ = protocol.Close() })

	openErr := make(chan error, 1)
	go func() {
		_, _, err := protocol.Open(context.Background())
		openErr <- err
	}()
	time.Sleep(20 * time.Millisecond)

	const config = uint32(1<<26) | uint32(1<<24) | 0x1F
	channel := []byte{0x2C, 0x01}
	send := func(timestamp int64, variable uint32) {
		t.Helper()
		transport.frames <- scriptFrame{data: buildV3Frame(timestamp, config, variable, channel)}
	}
	transport.frames <- scriptFrame{timeout: true}
	transport.frames <- scriptFrame{}
	for size := 1; size <= 20; size++ {
		transport.frames <- scriptFrame{data: make([]byte, size)}
	}
	send(0, 0x012F81AC)
	if !waitUntil(2*time.Second, func() bool { return protocol.GetDeviceId() == "012F81AC" }) {
		t.Fatalf("device id = %s", protocol.GetDeviceId())
	}
	send(1000, math.Float32bits(25.5))
	send(2000, math.Float32bits(121.5))
	send(3000, math.Float32bits(10.5))
	send(4000, math.Float32bits(18.5))
	if !waitUntil(2*time.Second, protocol.variablesReady) {
		t.Fatal("variables did not settle")
	}
	if got, err := protocol.GetTemperature(); err != nil || got != float64(float32(18.5)) {
		t.Fatalf("temperature = %v, %v", got, err)
	}
	send(5000, 0)
	select {
	case err := <-openErr:
		if err != nil {
			t.Fatalf("Open() error = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Open() did not become ready")
	}

	archived := make(chan Event, 2)
	realtime := make(chan Event, 16)
	if err := protocol.Subscribe("archive", func(event Event) { archived <- event }); err != nil {
		t.Fatal(err)
	}
	if err := protocol.SubscribeRealtime("live", func(event Event) { realtime <- event }); err != nil {
		t.Fatal(err)
	}
	for timestamp := int64(6000); timestamp < 16000; timestamp += 1000 {
		send(timestamp, 0)
	}
	select {
	case event := <-realtime:
		if len(event.ChannelData) != 1 || event.ChannelData[0].ChannelCode != "HNZ" || event.ChannelData[0].Data[0] != 300 {
			t.Fatalf("realtime event = %#v", event)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for realtime data")
	}
	select {
	case event := <-archived:
		if event.SampleRate != 10 || len(event.ChannelData[0].Data) != 10 {
			t.Fatalf("archived event = %#v", event)
		}
	case <-time.After(2 * time.Second):
		config := protocol.GetConfig()
		t.Fatalf("timed out waiting for archived data, rate=%d", config.GetSampleRate())
	}

	send(17000, 0)
	send(17100, 0)
	status := protocol.GetStatus()
	previousErrors := status.GetErrors()
	bad := buildV3Frame(18000, config, 0, channel)
	bad[0] ^= 0xFF
	transport.frames <- scriptFrame{data: bad}
	if !waitUntil(time.Second, func() bool {
		status := protocol.GetStatus()
		return status.GetErrors() > previousErrors
	}) {
		t.Fatal("corrupt packet did not increment errors")
	}
	if err := protocol.Close(); err != nil {
		t.Fatal(err)
	}
	assertStreamWorkersStopped(t, protocol.streamDone, protocol.decoderDone, protocol.ntpDone)
}

func TestV3OpenReadErrorCancels(t *testing.T) {
	transport := newScriptTransport()
	protocol := &ExplorerProtoImplV3{
		Logger:     testLogger(nil),
		TimeSource: timesource.New(time.Now),
		Transport:  transport,
		NtpOptions: NtpOptions{Pool: []string{startTestNTP(t)}, Retry: 0, ReadTimeout: 1},
	}
	t.Cleanup(func() { _ = protocol.Close() })
	openErr := make(chan error, 1)
	go func() {
		_, _, err := protocol.Open(context.Background())
		openErr <- err
	}()
	transport.frames <- scriptFrame{err: errors.New("read failed")}
	select {
	case err := <-openErr:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Open() error = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("read error did not cancel startup")
	}
	assertStreamWorkersStopped(t, protocol.streamDone, protocol.decoderDone, protocol.ntpDone)
}

func TestV3OpenNTPSync(t *testing.T) {
	transport := newScriptTransport()
	base := time.Unix(1_700_000_000, 0)
	protocol := &ExplorerProtoImplV3{
		ChannelCodes:    []string{"HNZ"},
		ExplorerOptions: ExplorerOptions{Protocol: "v3", Model: "E-C111G", Latitude: 25.5, Longitude: 121.5, Elevation: 10.5},
		NtpOptions:      NtpOptions{Pool: []string{startTestNTP(t)}, Retry: 0, ReadTimeout: 1},
		Logger:          testLogger(nil),
		TimeSource:      timesource.New(func() time.Time { return base }),
		Transport:       transport,
	}
	t.Cleanup(func() { _ = protocol.Close() })
	openErr := make(chan error, 1)
	go func() {
		_, _, err := protocol.Open(context.Background())
		openErr <- err
	}()

	const config = uint32(1<<24) | 0x1F
	channel := []byte{0x01, 0x00}
	const start = int64(1_600_000_000_000)
	send := func(timestamp int64, variable uint32) {
		t.Helper()
		transport.frames <- scriptFrame{data: buildV3Frame(timestamp, config, variable, channel)}
	}
	for i := 0; i < 10; i++ {
		send(start, 0x012F81AC)
	}
	if !waitUntil(2*time.Second, func() bool { return protocol.GetDeviceId() == "012F81AC" }) {
		t.Fatalf("device id %s", protocol.GetDeviceId())
	}
	send(start+1000, math.Float32bits(1))
	send(start+2000, math.Float32bits(2))
	send(start+3000, math.Float32bits(3))
	send(start+4000, math.Float32bits(4))
	if !waitUntil(2*time.Second, protocol.variablesReady) {
		t.Fatal("variables did not settle")
	}
	send(start+5000, 0)
	select {
	case err := <-openErr:
		if err != nil {
			t.Fatalf("Open() error = %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("NTP synchronization did not finish")
	}
	send(start+20000, 0)
	if !waitUntil(time.Second, func() bool { return !protocol.variablesReady() }) {
		t.Fatal("time jump did not reset variables")
	}
}

func newV3TestStream(t *testing.T) *ExplorerProtoImplV3 {
	t.Helper()
	protocol := &ExplorerProtoImplV3{
		Logger:             testLogger(nil),
		TimeSource:         timesource.New(time.Now),
		ChannelCodes:       []string{"HNZ"},
		fifoBuffer:         fifo.New[*explorerProtocolPacketV3](512),
		messageBus:         message.NewBus[Event](EXPLORER_STREAM_TOPIC),
		messageBusRealtime: message.NewBus[Event](EXPLORER_REALTIME_STREAM_TOPIC),
		stream:             explorerStreamV3{timeSourceInitialized: true},
		isDataStreamStable: true,
		variableAllSet:     true,
	}
	protocol.deviceConfig.SetGnssAvailability(true)
	id := uint32(0x012F81AC)
	protocol.deviceVariable.SetDeviceId(&id)
	t.Cleanup(protocol.messageBus.Close)
	t.Cleanup(protocol.messageBusRealtime.Close)
	return protocol
}

func v3PacketBody(timestamp int64, config uint32, channel []byte) []byte {
	frame := buildV3Frame(timestamp, config, 0, channel)
	return frame[:len(frame)-4]
}

func TestV3MalformedPacketsDoNotChangeState(t *testing.T) {
	const config = uint32(1<<26 | 1<<24)
	valid := v3PacketBody(1000, config, []byte{1, 0})
	corrupt := append([]byte(nil), valid...)
	corrupt[len(corrupt)-1] ^= 0xFF
	for _, tc := range []struct {
		name   string
		packet []byte
	}{
		{"empty", nil},
		{"short timestamp", []byte{1, 2, 3}},
		{"short config", make([]byte, 10)},
		{"short channels", v3PacketBody(1000, config, nil)},
		{"extra channels", v3PacketBody(1000, config, []byte{1, 0, 2, 0})},
		{"checksum", corrupt},
	} {
		t.Run(tc.name, func(t *testing.T) {
			protocol := newV3TestStream(t)
			protocol.prevMcuTimestamp = 100
			initFlag := int32(1)
			ready := make(chan struct{})
			now := timesource.MonotonicNow()
			protocol.processPacket(context.Background(), tc.packet, now, time.Now(), ready, &initFlag)
			protocol.decodePacket(&initFlag)
			status := protocol.GetStatus()
			if status.GetErrors() != 1 || status.GetFrames() != 0 || protocol.prevMcuTimestamp != 100 || protocol.GetDeviceId() != "012F81AC" || !protocol.variablesReady() {
				t.Fatal("invalid packet changed stream state or was not counted as an error")
			}
		})
	}
}

func TestV3TimeJumpDiscardsQueuedAndPartialArchive(t *testing.T) {
	protocol := newV3TestStream(t)
	archive := make(chan Event, 2)
	if err := protocol.Subscribe("archive", func(event Event) { archive <- event }); err != nil {
		t.Fatal(err)
	}
	const config = uint32(1<<26 | 1<<24) // 10 Hz, one int16 sample every 100 ms.
	initFlag := int32(1)
	ready := make(chan struct{})
	push := func(timestamp int64, sample byte, decode bool) {
		protocol.processPacket(context.Background(), v3PacketBody(timestamp, config, []byte{sample, 0}), timesource.MonotonicNow(), time.Now(), ready, &initFlag)
		if decode {
			protocol.decodePacket(&initFlag)
		}
	}
	for i := int64(0); i < 4; i++ {
		push(10_000+i*100, 7, true)
	}
	push(10_400, 7, false)
	push(10_500, 7, false)
	push(20_000, 9, true)
	if protocol.collectedSamples != 0 || len(protocol.channelDataBuf) != 0 || protocol.variablesReady() {
		t.Fatal("time jump left queued packets or a partial archive behind")
	}
	for i := int64(1); i <= 10; i++ {
		push(20_000+i*100, 9, true)
	}
	select {
	case event := <-archive:
		if event.Timestamp.UnixMilli() != 20_100 || !event.GNSSAvailable || event.SampleRate != 10 || len(event.ChannelData) != 1 || len(event.ChannelData[0].Data) != 10 {
			t.Fatalf("archive after reset = %#v", event)
		}
		for _, sample := range event.ChannelData[0].Data {
			if sample != 9 {
				t.Fatal("archive contains samples from before the time jump")
			}
		}
	case <-time.After(time.Second):
		t.Fatal("stream did not recover after the time jump")
	}
}

func TestV3ConcurrentReaderResetAndDecoder(t *testing.T) {
	protocol := newV3TestStream(t)
	const config = uint32(1<<26 | 1<<24)
	initFlag := int32(1)
	ready := make(chan struct{})
	var workers sync.WaitGroup
	workers.Add(2)
	go func() {
		defer workers.Done()
		for i := int64(0); i < 1000; i++ {
			// Periodically jump far enough to reset a concurrently decoded archive.
			timestamp := 10_000 + i*100 + (i/20)*10_000
			protocol.processPacket(context.Background(), v3PacketBody(timestamp, config, []byte{1, 0}), timesource.MonotonicNow(), time.Now(), ready, &initFlag)
		}
	}()
	go func() {
		defer workers.Done()
		for i := 0; i < 1000; i++ {
			protocol.decodePacket(&initFlag)
			_ = protocol.variablesReady()
			_ = protocol.GetConfig()
			_ = protocol.GetStatus()
		}
	}()
	workers.Wait()
	protocol.processPacket(context.Background(), v3PacketBody(1_000_000, config, []byte{2, 0}), timesource.MonotonicNow(), time.Now(), ready, &initFlag)
	protocol.decodePacket(&initFlag)
	if protocol.collectedSamples != 0 || len(protocol.channelDataBuf) != 0 || protocol.variablesReady() {
		t.Fatal("decoder repopulated a reset archive")
	}
}
