package explorer

import (
	"context"
	"encoding/binary"
	"testing"
	"testing/synctest"
	"time"

	"github.com/anyshake/observer/pkg/fifo"
	"github.com/anyshake/observer/pkg/message"
)

// receiveJitter returns a deterministic, zero-mean serial read jitter of up to
// ±7 ms that changes every packet, like the timestamps in the NTP-mode report.
func receiveJitter(i int) time.Duration {
	return time.Duration((i*7919)%15-7) * time.Millisecond
}

// assertContinuous checks that each archive record starts exactly one second
// after the previous one, within the bounded slew and millisecond rounding.
// Overlapping records make SeedLink drop samples; the old code failed this.
func assertContinuous(t *testing.T, starts []time.Time) {
	t.Helper()
	if len(starts) < 3 {
		t.Fatalf("only %d archive records", len(starts))
	}
	for i := 1; i < len(starts); i++ {
		if d := starts[i].Sub(starts[i-1]) - time.Second; d < -time.Millisecond || d > time.Millisecond {
			t.Fatalf("archive %d starts %s after the previous record ends", i, d)
		}
	}
}

// collectArchive waits for the asynchronous message bus to deliver every
// published archive record and returns their start times in order.
func collectArchive(t *testing.T, archive <-chan Event, published int64) []time.Time {
	t.Helper()
	starts := make([]time.Time, 0, published)
	for int64(len(starts)) < published {
		select {
		case event := <-archive:
			starts = append(starts, event.Timestamp)
		case <-time.After(2 * time.Second):
			t.Fatalf("received %d of %d archive records", len(starts), published)
		}
	}
	return starts
}

func TestV1NTPArchiveIsContinuousUnderJitter(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		protocol := &ExplorerProtoImplV1{
			Logger:             testLogger(nil),
			ChannelCodes:       []string{"HNZ", "HNE", "HNN"},
			fifoBuffer:         fifo.New[byte](74 * 32),
			decoderDone:        make(chan struct{}),
			messageBus:         message.NewBus[Event](EXPLORER_STREAM_TOPIC),
			messageBusRealtime: message.NewBus[Event](EXPLORER_REALTIME_STREAM_TOPIC),
		}
		protocol.deviceConfig.SetSampleRate(50)
		protocol.deviceConfig.SetPacketInterval(100 * time.Millisecond)
		defer protocol.messageBus.Close()
		defer protocol.messageBusRealtime.Close()
		go protocol.decodeStream(ctx, 66)
		defer func() {
			cancel()
			<-protocol.decoderDone
		}()
		archive := make(chan Event, 64)
		if err := protocol.Subscribe("archive", func(event Event) { archive <- event }); err != nil {
			t.Fatal(err)
		}
		const start = int64(1_800_000_000_000)
		for i := 0; i < 300; i++ {
			samples := make([]int32, 5)
			packet := buildV1DataPacket(2, 5, samples, samples, samples)
			timestamp := start + int64(i)*100 + receiveJitter(i).Milliseconds()
			binary.BigEndian.PutUint64(packet[2:10], uint64(timestamp))
			_, _ = protocol.fifoBuffer.Write(packet...)
			time.Sleep(10 * time.Millisecond)
			synctest.Wait()
		}
		var starts []time.Time
		for len(archive) > 0 {
			starts = append(starts, (<-archive).Timestamp)
		}
		assertContinuous(t, starts)
	})
}

func TestV2NTPArchiveIsContinuousAcrossOffsetStep(t *testing.T) {
	s := newV2TestStream(t)
	s.startNTP(100_000, 1)
	base := s.received
	// Discard records published while the stream was starting.
	before := s.protocol.deviceStatus.GetMessages()
	collectArchive(t, s.archive, before)
	for i := 0; i < 400; i++ {
		// Stable receive timing lets v2 refresh its averaged offset. The
		// measured time steps 8 ms at packet 150, as after an NTP update.
		at := base.Add(time.Duration(i)*100*time.Millisecond + time.Duration(i%3)*time.Millisecond)
		if i >= 150 {
			at = at.Add(8 * time.Millisecond)
		}
		s.received = at
		s.push(101_100+int64(i)*100, false, 1)
	}
	assertContinuous(t, collectArchive(t, s.archive, s.protocol.deviceStatus.GetMessages()-before))
}

func TestV3NTPArchiveIsContinuousUnderJitter(t *testing.T) {
	protocol := newV3TestStream(t)
	protocol.deviceConfig.SetGnssAvailability(false)
	archive := make(chan Event, 64)
	if err := protocol.Subscribe("archive", func(event Event) { archive <- event }); err != nil {
		t.Fatal(err)
	}
	// 250 Hz, 100 ms packets, one int16 channel, no measured delay.
	const config = uint32(5<<27 | 1<<24)
	channel := make([]byte, 25*2)
	initFlag := int32(1)
	received := time.Unix(1_800_000_000, 0)
	for i := 0; i < 300; i++ {
		frame := buildV3Frame(int64(1000+i*100), config, 0, channel)
		protocol.processPacket(context.Background(), frame[:len(frame)-4], received, received.Add(time.Duration(i)*100*time.Millisecond+receiveJitter(i)), make(chan struct{}), &initFlag)
		protocol.decodePacket(&initFlag)
	}
	assertContinuous(t, collectArchive(t, archive, protocol.deviceStatus.GetMessages()))
}
