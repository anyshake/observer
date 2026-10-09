package explorer

import (
	"bytes"
	"context"
	"math"
	"testing"
	"time"

	"github.com/anyshake/observer/pkg/message"
	"github.com/anyshake/observer/pkg/ringbuf"
	"github.com/anyshake/observer/pkg/timesource"
)

const v2TestEpoch = int64(1_700_000_000_000)

type v2TestStream struct {
	t        *testing.T
	protocol *ExplorerProtoImplV2
	live     chan Event
	archive  chan Event
	received time.Time
}

func newV2TestStream(t *testing.T) *v2TestStream {
	t.Helper()
	now := timesource.MonotonicNow()
	g := &ExplorerProtoImplV2{
		Logger:             testLogger(nil),
		TimeSource:         timesource.New(time.Now),
		ChannelCodes:       []string{"HNZ", "HNE", "HNN"},
		ExplorerOptions:    ExplorerOptions{Latitude: 25, Longitude: 121, Elevation: 10},
		messageBus:         message.NewBus[Event](EXPLORER_STREAM_TOPIC),
		messageBusRealtime: message.NewBus[Event](EXPLORER_REALTIME_STREAM_TOPIC),
		clockDriftBuf:      ringbuf.New[clockDrift](100),
		// Inject an already measured NTP offset; packet tests need no network.
		stream: explorerStreamV2{ntpMeasuredAt: now, ntpOffset: time.Now().Sub(now)},
	}
	t.Cleanup(g.messageBus.Close)
	t.Cleanup(g.messageBusRealtime.Close)
	s := &v2TestStream{t: t, protocol: g, live: make(chan Event, 128), archive: make(chan Event, 128), received: now}
	if err := g.SubscribeRealtime("test", func(event Event) { s.live <- event }); err != nil {
		t.Fatal(err)
	}
	if err := g.Subscribe("test", func(event Event) { s.archive <- event }); err != nil {
		t.Fatal(err)
	}
	return s
}

func v2TestVariable(timestamp int64, gnss bool) uint32 {
	switch (timestamp / 1000) % 4 {
	case 0:
		id := uint32(0x012F81AC)
		if gnss {
			id |= 0x80000000
		}
		return id
	case 1:
		return math.Float32bits(25.5)
	case 2:
		return math.Float32bits(121.5)
	default:
		return math.Float32bits(10.5)
	}
}

func (s *v2TestStream) push(timestamp int64, gnss bool, sample int32) *Event {
	s.t.Helper()
	packet := buildV2WirePacket(timestamp, v2TestVariable(timestamp, gnss), sample)
	ready, err := s.protocol.processPacket(context.Background(), packet, s.received, 0)
	s.received = s.received.Add(100 * time.Millisecond)
	if err != nil {
		s.t.Fatal(err)
	}
	if !ready {
		return nil
	}
	select {
	case event := <-s.live:
		return &event
	case <-time.After(time.Second):
		s.t.Fatal("missing realtime event after processing a valid packet")
		return nil
	}
}

func (s *v2TestStream) startNTP(timestamp int64, sample int32) *Event {
	s.t.Helper()
	for i := int64(0); i < 10; i++ {
		s.push(timestamp+i*100, false, sample)
	}
	event := s.push(timestamp+1000, false, sample)
	if event == nil {
		s.t.Fatal("NTP stream did not initialize after a full second")
	}
	return event
}

func TestV2NTPToGNSSWaitsForDeviceWord(t *testing.T) {
	s := newV2TestStream(t)
	s.startNTP(100_000, 7)
	s.push(101_100, false, 7) // Leave an incomplete NTP archive.
	for _, timestamp := range []int64{v2TestEpoch + 1000, v2TestEpoch + 2000, v2TestEpoch + 3000} {
		if event := s.push(timestamp, true, 9); event != nil {
			t.Fatalf("published unconfirmed GNSS timestamp with NTP offset: %+v", event)
		}
	}
	for timestamp := v2TestEpoch + 4000; timestamp < v2TestEpoch+7000; timestamp += 1000 {
		s.push(timestamp, true, 9)
	}
	event := s.push(v2TestEpoch+7000, true, 9)
	if event == nil || !event.GNSSAvailable || event.Timestamp.UnixMilli() != v2TestEpoch+7000 {
		t.Fatalf("GNSS event = %+v", event)
	}
	select {
	case archived := <-s.archive:
		if archived.Timestamp != event.Timestamp || archived.SampleRate != 5 || !archived.GNSSAvailable {
			t.Fatalf("archive crossed the mode boundary: %+v", archived)
		}
		for _, channel := range archived.ChannelData {
			if len(channel.Data) != 5 {
				t.Fatalf("mixed archive samples: %v", channel.Data)
			}
			for _, sample := range channel.Data {
				if sample != 9 {
					t.Fatalf("old-mode sample retained: %v", channel.Data)
				}
			}
		}
	case <-time.After(time.Second):
		t.Fatal("missing GNSS archive")
	}
}

func TestV2GNSSToNTPAndMCUReboot(t *testing.T) {
	for _, initialGNSS := range []bool{false, true} {
		t.Run(map[bool]string{false: "MCU reboot", true: "GNSS to NTP"}[initialGNSS], func(t *testing.T) {
			s := newV2TestStream(t)
			start := int64(100_000)
			if initialGNSS {
				start = v2TestEpoch
			}
			for i := int64(0); i <= 3; i++ {
				s.push(start+i*1000, initialGNSS, 1)
			}
			for timestamp := int64(1000); timestamp < 4000; timestamp += 1000 {
				if event := s.push(timestamp, false, 2); event != nil {
					t.Fatalf("published uptime before new DEVICE_ID: %+v", event)
				}
			}
			s.push(4000, false, 2)
			for timestamp := int64(4100); timestamp < 5000; timestamp += 100 {
				s.push(timestamp, false, 2)
			}
			event := s.push(5000, false, 2)
			if event == nil || event.GNSSAvailable {
				t.Fatalf("NTP recovery event = %+v", event)
			}
			if distance := time.Since(event.Timestamp); distance < -time.Minute || distance > time.Minute {
				t.Fatalf("old clock/offset survived recovery: %v", event.Timestamp)
			}
		})
	}
}

func TestV2RejectsCorruptPacketBeforeUpdatingState(t *testing.T) {
	s := newV2TestStream(t)
	s.startNTP(100_000, 1)
	before := s.protocol.TimeSource.Now()
	generation := s.protocol.stream.generation
	packet := buildV2WirePacket(v2TestEpoch, 0x80000055, 2)
	packet[len(packet)-1] ^= 1
	if _, err := s.protocol.processPacket(context.Background(), packet, s.received, 0); err == nil {
		t.Fatal("corrupt checksum accepted")
	}
	if s.protocol.stream.generation != generation || s.protocol.deviceConfig.GetGnssAvailability() || s.protocol.GetDeviceId() != "012F81AC" {
		t.Fatal("corrupt packet altered the time domain or device ID")
	}
	if distance := s.protocol.TimeSource.Now().Sub(before); distance < 0 || distance > time.Second {
		t.Fatal("corrupt packet adjusted the time source")
	}
	if event := s.push(101_100, false, 1); event == nil {
		t.Fatal("corrupt packet interrupted the valid stream")
	}
}

func TestV2RejectsUnrepresentableTimestamp(t *testing.T) {
	for _, timestamp := range []int64{-1, math.MaxInt64, 1 << 60} {
		s := newV2TestStream(t)
		s.startNTP(100_000, 1)
		if event := s.push(timestamp, true, 2); event != nil {
			t.Fatalf("published invalid timestamp %d: %+v", timestamp, event)
		}
		if s.protocol.variablesReady() || s.protocol.stream.modeKnown {
			t.Fatal("invalid timestamp did not invalidate the time domain")
		}
	}
}

func TestV2ModeChangeWithoutTimestampJump(t *testing.T) {
	s := newV2TestStream(t)
	// Even if timestamps happen to remain continuous, the DEVICE_ID mode bit
	// must invalidate the old mapping and previously collected coordinates.
	s.startNTP(v2TestEpoch+4000, 1)
	s.push(v2TestEpoch+7900, false, 1)
	if event := s.push(v2TestEpoch+8000, true, 2); event != nil {
		t.Fatal("mode change immediately reused previous time/metadata")
	}
	if s.protocol.variablesReady() || s.protocol.stream.clockInitialized || len(s.protocol.channelDataBuf) != 0 {
		t.Fatal("mode change did not clear stale state")
	}
}

func TestV2LateNTPMeasurementCannotOverwriteNewMode(t *testing.T) {
	s := newV2TestStream(t)
	s.startNTP(100_000, 1)
	oldGeneration := s.protocol.stream.generation
	for timestamp := v2TestEpoch; timestamp <= v2TestEpoch+3000; timestamp += 1000 {
		s.push(timestamp, true, 2)
	}
	before := s.protocol.TimeSource.Now()
	if s.protocol.applyNTPMeasurement(oldGeneration, time.Hour) {
		t.Fatal("old NTP measurement replaced the GNSS clock")
	}
	if distance := s.protocol.TimeSource.Now().Sub(before); distance < 0 || distance > time.Second {
		t.Fatal("rejected NTP measurement altered the clock")
	}
	s.startNTP(4000, 3)
	if s.protocol.applyNTPMeasurement(oldGeneration, time.Hour) {
		t.Fatal("old NTP measurement survived a round-trip mode switch")
	}
	if !s.protocol.applyNTPMeasurement(s.protocol.stream.generation, s.protocol.stream.ntpOffset) {
		t.Fatal("current NTP measurement was rejected")
	}
}

func TestV2SampleRateToleratesTimestampJitter(t *testing.T) {
	s := newV2TestStream(t)
	s.push(100_000, false, 1)
	for i := int64(1); i <= 50; i++ {
		s.push(100_000+i*40+i%2, false, 1)
	}
	select {
	case event := <-s.archive:
		if event.SampleRate != 125 || len(event.ChannelData[0].Data) != 125 {
			t.Fatalf("jitter changed the sample rate or archive size: %+v", event)
		}
	case <-time.After(time.Second):
		t.Fatal("jitter prevented one-second archives")
	}
}

func TestV2DuplicateAndMissingPacketsDoNotEnterArchive(t *testing.T) {
	s := newV2TestStream(t)
	s.startNTP(100_000, 1)
	if event := s.push(101_000, false, 2); event != nil {
		t.Fatal("duplicate packet was published")
	}
	if event := s.push(101_300, false, 2); event != nil {
		t.Fatal("packet after a gap was appended to the old archive")
	}
	for timestamp := int64(101_400); timestamp <= 103_200; timestamp += 100 {
		s.push(timestamp, false, 3)
	}
	select {
	case event := <-s.archive:
		if event.SampleRate != 50 || len(event.ChannelData[0].Data) != 50 {
			t.Fatalf("archive after gap = %+v", event)
		}
		for _, value := range event.ChannelData[0].Data {
			if value != 3 {
				t.Fatal("archive contains pre-gap samples")
			}
		}
	case <-time.After(time.Second):
		t.Fatal("stream did not recover after gap")
	}
}

func TestV2ReadStreamHandlesSplitAndCoalescedPackets(t *testing.T) {
	transport := newScriptTransport()
	g := &ExplorerProtoImplV2{
		Transport:  transport,
		Logger:     testLogger(nil),
		TimeSource: timesource.New(timesource.MonotonicNow),
		NtpOptions: NtpOptions{Pool: []string{startTestNTP(t)}, ReadTimeout: 1},
	}
	ctx, cancel := context.WithCancel(context.Background())
	opened := make(chan error, 1)
	go func() {
		_, _, err := g.Open(ctx)
		opened <- err
	}()
	t.Cleanup(func() {
		cancel()
		_ = transport.Close()
		<-opened
		_ = g.Close()
	})
	var wire []byte
	for i := int64(0); i <= 5; i++ {
		timestamp := v2TestEpoch + i*1000
		wire = append(wire, buildV2WirePacket(timestamp, v2TestVariable(timestamp, true), int32(i))...)
	}
	// Include a false header and junk, then split the real header and payload.
	wire = append([]byte{0xFA, 0xDE, 0, 0, 0}, wire...)
	transport.packets <- wire[:6]
	transport.packets <- wire[6:37]
	for rest := wire[37:]; len(rest) > 0; {
		n := min(len(rest), V2_PACKET_SIZE*2)
		transport.packets <- bytes.Clone(rest[:n])
		rest = rest[n:]
	}
	select {
	case err := <-opened:
		opened <- err // Cleanup joins Open even when startup fails.
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("split/coalesced stream did not initialize")
	}
	if !waitUntil(time.Second, func() bool {
		status := g.GetStatus()
		return status.GetUpdatedAt().UnixMilli() == v2TestEpoch+5000
	}) {
		t.Fatal("not all packets in the transport stream were decoded")
	}
}
