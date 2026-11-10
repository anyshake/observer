package explorer

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/anyshake/observer/pkg/timesource"
)

func packetDelayWord(delay int32) [4]byte {
	var word [4]byte
	binary.LittleEndian.PutUint32(word[:], uint32(delay))
	return word
}

type v3TimingTransport struct{ fakeTransport }

func (*v3TimingTransport) GetLatency(bytes int) time.Duration {
	return time.Duration(bytes) * 87 * time.Microsecond
}

func TestV3PacketDelayMetadata(t *testing.T) {
	const config = uint32(5<<27 | 1<<24 | 1<<5)
	protocol := &ExplorerProtoImplV3{}
	collect := func(timestamp int64, config uint32, word [4]byte) {
		protocol.updatePacketDelay(timestamp, config, binary.LittleEndian.Uint32(word[:]))
		protocol.getVariableData(timestamp, config, word)
	}
	collect(4000, config, [4]byte{})
	if protocol.variablesReady() {
		t.Fatal("metadata settled before a delay measurement arrived")
	}
	collect(5000, config, packetDelayWord(207123))
	if !protocol.variablesReady() || protocol.stream.packetDelay != 207123*time.Microsecond {
		t.Fatal("valid microsecond delay was not collected")
	}
	collect(6000, config, [4]byte{})
	if !protocol.variablesReady() || protocol.stream.packetDelay != 207123*time.Microsecond {
		t.Fatal("reserved metadata slot changed the cached delay")
	}
	for _, delay := range []int32{-1, math.MaxInt32} {
		collect(5000, config, packetDelayWord(delay))
		if protocol.variablesReady() {
			t.Fatalf("unavailable delay %d was accepted", delay)
		}
	}
	collect(5000, config, packetDelayWord(0))
	if !protocol.variablesReady() || protocol.stream.packetDelay != 0 {
		t.Fatal("a measured zero delay must be valid")
	}
	protocol.resetVariables()
	collect(6000, config, [4]byte{})
	if protocol.variablesReady() || protocol.stream.packetDelayKnown {
		t.Fatal("device reset retained a delay measurement")
	}

	for _, change := range []uint32{1 << 27, 1 << 30, 1 << 22, 1 << 26} {
		collect(5000, config, packetDelayWord(207123))
		collect(6000, config^change, [4]byte{})
		if protocol.variablesReady() || protocol.stream.packetDelayKnown {
			t.Fatalf("configuration change %#x retained a stale delay", change)
		}
	}
	collect(5000, config, packetDelayWord(207123))
	collect(6000, config&^(1<<5), [4]byte{})
	if !protocol.variablesReady() || protocol.stream.packetDelayKnown {
		t.Fatal("legacy firmware must not require a delay field")
	}
}

func TestV3MeasuredDelaySampleTimestamps(t *testing.T) {
	for rateIndex, rate := range map[uint32]int{2: 50, 3: 100, 4: 200, 5: 250} {
		for _, gnss := range []bool{false, true} {
			for _, measured := range []bool{false, true} {
				name := fmt.Sprintf("%dHz/GNSS=%v/measured=%v", rate, gnss, measured)
				t.Run(name, func(t *testing.T) {
					protocol := newV3TestStream(t)
					protocol.Transport = &v3TimingTransport{}
					config := rateIndex<<27 | 1<<24
					if gnss {
						config |= 1 << 26
					}
					if measured {
						config |= 1 << 5
					}
					// Includes a nonintegral millisecond and the additional 80 ms
					// FIR delay only for the rational resampling rates.
					delay := 207123 * time.Microsecond
					if rate == 100 || rate == 200 {
						delay += 80 * time.Millisecond
					}
					const timestamp = int64(15000) // Slot 5.
					channel := make([]byte, rate/10*2)
					for i := 0; i < len(channel); i += 2 {
						binary.LittleEndian.PutUint16(channel[i:], 123)
					}
					frame := buildV3Frame(timestamp, config, uint32(delay.Microseconds()), channel)
					received := time.Unix(1800000000, 987654321)
					initFlag := int32(1)
					events := make(chan Event, 1)
					if err := protocol.SubscribeRealtime("timing", func(event Event) { events <- event }); err != nil {
						t.Fatal(err)
					}
					protocol.processPacket(context.Background(), frame[:len(frame)-4], timesource.MonotonicNow(), received, make(chan struct{}), &initFlag)
					protocol.decodePacket(&initFlag)
					want := received.Add(-200 * time.Millisecond)
					if measured {
						want = received.Add(-100*time.Millisecond - delay - 174*time.Microsecond)
					}
					if gnss {
						want = time.UnixMilli(timestamp)
					}
					select {
					case event := <-events:
						if !event.Timestamp.Equal(want) {
							t.Fatalf("timestamp = %v, want %v", event.Timestamp, want)
						}
						if event.SampleRate != rate || len(event.ChannelData) != 1 || len(event.ChannelData[0].Data) != rate/10 || event.ChannelData[0].Data[0] != 123 {
							t.Fatalf("timing correction changed sample data: %#v", event)
						}
					case <-time.After(time.Second):
						t.Fatal("no sample event received")
					}
				})
			}
		}
	}
}

func TestV3WaitingForDelayDiscardsPartialArchive(t *testing.T) {
	protocol := newV3TestStream(t)
	protocol.collectedSamples = 10
	protocol.packetTimeObj = time.UnixMilli(1000)
	protocol.channelDataBuf = []ChannelData{{Data: make([]int32, 10)}}
	const config = uint32(5<<27 | 1<<24 | 1<<5)
	initFlag := int32(1)
	protocol.processPacket(context.Background(), v3PacketBody(6000, config, make([]byte, 50)), timesource.MonotonicNow(), time.Now(), make(chan struct{}), &initFlag)
	protocol.decodePacket(&initFlag)
	if protocol.variablesReady() || protocol.collectedSamples != 0 || !protocol.packetTimeObj.IsZero() || len(protocol.channelDataBuf) != 0 {
		t.Fatal("waiting for timing metadata retained an archive across a gap")
	}
	status := protocol.GetStatus()
	if status.GetFrames() != 0 {
		t.Fatal("a packet without a known delay was published")
	}
}

func TestV3GNSSClockUsesMeasuredDelay(t *testing.T) {
	for rateIndex, rate := range map[uint32]int{2: 50, 3: 100, 4: 200, 5: 250} {
		for _, measured := range []bool{false, true} {
			t.Run(fmt.Sprintf("%dHz/measured=%v", rate, measured), func(t *testing.T) {
				protocol := newV3TestStream(t)
				protocol.Transport = &v3TimingTransport{}
				protocol.stream.timeSourceInitialized = false
				protocol.timeCalibrationChan = make(chan [2]time.Time, 1)
				config := rateIndex<<27 | 1<<26 | 1<<24
				delay := 207123 * time.Microsecond
				if rate == 100 || rate == 200 {
					delay += 80 * time.Millisecond
				}
				latency := 200 * time.Millisecond
				if measured {
					config |= 1 << 5
					latency = 100*time.Millisecond + delay + 174*time.Microsecond
				}
				push := func(timestamp int64) time.Time {
					frame := buildV3Frame(timestamp, config, uint32(delay.Microseconds()), make([]byte, rate/10*2))
					received := timesource.MonotonicNow()
					initFlag := int32(1)
					protocol.processPacket(context.Background(), frame[:len(frame)-4], received, time.Now(), make(chan struct{}), &initFlag)
					return received
				}
				const timestamp = int64(15000)
				received := push(timestamp)
				if !protocol.stream.timeSourceInitialized {
					t.Fatal("GNSS clock was not initialized")
				}
				want := time.UnixMilli(timestamp).Add(latency)
				before := timesource.MonotonicNow()
				now := protocol.TimeSource.Now()
				after := timesource.MonotonicNow()
				if now.Before(want.Add(before.Sub(received))) || now.After(want.Add(after.Sub(received))) {
					t.Fatalf("shared clock = %v, want receive-time UTC near %v", now, want)
				}
				// Subsequent calibrations must use the same full latency as
				// initialization, including submillisecond header transit.
				received = push(timestamp + 100)
				select {
				case pair := <-protocol.timeCalibrationChan:
					if !pair[0].Equal(received) || !pair[1].Equal(time.UnixMilli(timestamp+100).Add(latency)) {
						t.Fatalf("periodic clock calibration = %v", pair)
					}
				default:
					t.Fatal("no periodic GNSS clock calibration")
				}
			})
		}
	}
}

func TestV3QueuedPacketWaitsForItsOwnTiming(t *testing.T) {
	protocol := newV3TestStream(t)
	protocol.Transport = &v3TimingTransport{}
	const config = uint32(5<<27 | 1<<24 | 1<<5)
	initFlag := int32(1)
	ready := make(chan struct{})
	received := time.Unix(1800000000, 0)
	protocol.processPacket(context.Background(), v3PacketBody(14900, config, make([]byte, 50)), timesource.MonotonicNow(), received, ready, &initFlag)
	frame := buildV3Frame(15000, config, 207123, make([]byte, 50))
	protocol.processPacket(context.Background(), frame[:len(frame)-4], timesource.MonotonicNow(), received.Add(100*time.Millisecond), ready, &initFlag)
	protocol.decodePacket(&initFlag)
	status := protocol.GetStatus()
	if status.GetFrames() != 0 {
		t.Fatal("a queued packet was timestamped before its delay was available")
	}
	protocol.decodePacket(&initFlag)
	status = protocol.GetStatus()
	want := received.Add(-207123*time.Microsecond - 174*time.Microsecond)
	if status.GetFrames() != 1 || !status.GetUpdatedAt().Equal(want) {
		t.Fatalf("timestamp after timing became available = %v, want %v", status.GetUpdatedAt(), want)
	}
}
