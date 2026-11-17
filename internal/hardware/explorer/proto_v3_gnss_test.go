package explorer

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/anyshake/observer/pkg/timesource"
)

func TestV3GnssStatusMetadata(t *testing.T) {
	t.Parallel()

	zeroHDOP, lockedHDOP, holdoverHDOP := 0.0, 0.9, 1.2
	zeroAge, lockedAge, holdoverAge := 0, 2, 300
	for _, tc := range []struct {
		name string
		word [4]byte
		want GnssStatus
	}{
		{"locked", [4]byte{0x4F, 0x60, 0x02, 0x00}, GnssStatus{Locked: true, FixValid: true, TimeValid: true, HDOP: &lockedHDOP, Satellites: 12, PPSAgeSeconds: &lockedAge}},
		{"holdover", [4]byte{0x60, 0x28, 0x2C, 0x01}, GnssStatus{HDOP: &holdoverHDOP, Satellites: 5, PPSAgeSeconds: &holdoverAge}},
		{"saturated", [4]byte{0xFF, 0xFF, 0xFF, 0xFF}, GnssStatus{Locked: true, FixValid: true, TimeValid: true, Satellites: 31}},
		{"zero", [4]byte{}, GnssStatus{HDOP: &zeroHDOP, PPSAgeSeconds: &zeroAge}},
		{"lock only", [4]byte{1}, GnssStatus{Locked: true, HDOP: &zeroHDOP, PPSAgeSeconds: &zeroAge}},
		{"fix only", [4]byte{2}, GnssStatus{FixValid: true, HDOP: &zeroHDOP, PPSAgeSeconds: &zeroAge}},
		{"time only", [4]byte{4}, GnssStatus{TimeValid: true, HDOP: &zeroHDOP, PPSAgeSeconds: &zeroAge}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			protocol := &ExplorerProtoImplV3{}
			const config = uint32(1<<26 | 1<<6)
			protocol.getVariableData(6000, config, tc.word)
			if got := protocol.GetStatus().GetGnssStatus(); !reflect.DeepEqual(got, &tc.want) {
				t.Fatalf("GNSS status = %#v, want %#v", got, &tc.want)
			}
			protocol.getVariableData(7000, config, [4]byte{0xFF, 0xFF, 0xFF, 0xFF})
			if got := protocol.GetStatus().GetGnssStatus(); !reflect.DeepEqual(got, &tc.want) {
				t.Fatal("reserved metadata slot changed GNSS status")
			}
		})
	}
}

func TestV3GnssStatusLifecycle(t *testing.T) {
	protocol := &ExplorerProtoImplV3{}
	const config = uint32(1<<26 | 1<<6)
	word := [4]byte{0x4F, 0x60, 0x02, 0x00}
	protocol.getVariableData(6000, config&^(1<<6), word)
	if protocol.GetStatus().GetGnssStatus() != nil {
		t.Fatal("legacy firmware exposed GNSS status without its capability bit")
	}
	protocol.getVariableData(5000, config, word)
	if protocol.GetStatus().GetGnssStatus() != nil || !protocol.variablesReady() {
		t.Fatal("GNSS diagnostics must remain unknown without blocking waveform readiness")
	}
	protocol.getVariableData(6000, config, word)
	snapshot := protocol.GetStatus()
	status := snapshot.GetGnssStatus()
	*status.HDOP, *status.PPSAgeSeconds = 99, 99
	if got := snapshot.GetGnssStatus(); *got.HDOP != 0.9 || *got.PPSAgeSeconds != 2 {
		t.Fatal("GNSS status getter exposed mutable state")
	}
	snapshot.SetGnssStatus(status)
	*status.HDOP, *status.PPSAgeSeconds = 100, 100
	if got := snapshot.GetGnssStatus(); *got.HDOP != 99 || *got.PPSAgeSeconds != 99 {
		t.Fatal("GNSS status setter retained caller-owned values")
	}
	if got := protocol.GetStatus().GetGnssStatus(); *got.HDOP != 0.9 || *got.PPSAgeSeconds != 2 {
		t.Fatal("device status snapshot changed the live GNSS status")
	}
	protocol.getVariableData(7000, config&^(1<<6), word)
	if protocol.GetStatus().GetGnssStatus() != nil {
		t.Fatal("disabled capability retained a stale GNSS status")
	}
	protocol.getVariableData(6000, config, word)
	protocol.resetVariables()
	if protocol.GetStatus().GetGnssStatus() != nil || snapshot.GetGnssStatus() == nil {
		t.Fatal("device reset did not clear live GNSS status independently of snapshots")
	}
}

func TestV3GnssHoldoverKeepsUTCTimestamps(t *testing.T) {
	protocol := newV3TestStream(t)
	const config = uint32(1<<26 | 1<<24 | 1<<6)
	const timestamp = int64(16000)
	// Holdover, no valid fix or NMEA time, HDOP 1.2, 5 satellites, PPS age 300s.
	frame := buildV3Frame(timestamp, config, 0x012C2860, []byte{123, 0})
	events := make(chan Event, 1)
	if err := protocol.SubscribeRealtime("holdover", func(event Event) { events <- event }); err != nil {
		t.Fatal(err)
	}
	initFlag := int32(1)
	protocol.processPacket(context.Background(), frame[:len(frame)-4], timesource.MonotonicNow(), time.Unix(1800000000, 0), make(chan struct{}), &initFlag)
	protocol.decodePacket(&initFlag)
	if status := protocol.GetStatus().GetGnssStatus(); status == nil || status.Locked || status.FixValid || status.TimeValid {
		t.Fatalf("holdover status = %#v", status)
	}
	select {
	case event := <-events:
		if !event.GNSSAvailable || !event.Timestamp.Equal(time.UnixMilli(timestamp)) || event.ChannelData[0].Data[0] != 123 {
			t.Fatalf("holdover changed UTC sample timestamps or waveform data: %#v", event)
		}
	case <-time.After(time.Second):
		t.Fatal("GNSS holdover stopped waveform publication")
	}
}
