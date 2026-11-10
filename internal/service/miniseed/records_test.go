package miniseed

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/waveform"
	"github.com/bclswl0827/mseedio"
)

func TestSaveMiniSeedRecordsWritesTimeCorrection(t *testing.T) {
	dir := t.TempDir()
	start := time.Date(2024, 3, 15, 13, 5, 9, 123000000, time.UTC)
	jitters := []time.Duration{0, 15 * time.Millisecond, 30 * time.Millisecond}
	svc := &MiniSeedServiceImpl{
		filePath:     dir,
		stationCode:  "SHAKE",
		networkCode:  "AS",
		locationCode: "00",
		dataSequence: sequence{
			filePath:     filepath.Join(dir, ".miniseed_sequence"),
			sequenceData: make(map[string]uint32),
		},
		timelines: make(map[string]*waveform.Timeline),
	}

	for i, jitter := range jitters {
		timestamp := start.Add(time.Duration(i)*time.Second + jitter)
		svc.recordBuffer = [][]buffer{{{
			Timestamp:  timestamp.UnixMilli(),
			SampleRate: 4,
			ChannelData: explorer.ChannelData{
				ChannelCode: "EHZ",
				Data:        []int32{1, -2, 3, -4},
			},
		}}}
		if _, err := svc.saveMiniSeedRecords(false); err != nil {
			t.Fatal(err)
		}
	}

	filePath := filepath.Join(dir, "2024-03-15", svc.getMiniSeedFileName(start, "EHZ"))
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	var decoded mseedio.MiniSeedData
	if err := decoded.ReadFromReader(bytes.NewReader(data)); err != nil {
		t.Fatal(err)
	}
	if len(decoded.Series) != 3 {
		t.Fatalf("series count = %d, want 3", len(decoded.Series))
	}
	for i, series := range decoded.Series {
		wantTime := start.Add(time.Duration(i) * time.Second)
		recordTime := wantTime.Add(jitters[i])
		if !series.FixedSection.StartTime.Equal(wantTime) {
			t.Fatalf("record %d effective start time = %s, want %s", i, series.FixedSection.StartTime, wantTime)
		}
		header := data[series.FixedSection.ReaderOffset.Start:]
		gotFraction := binary.BigEndian.Uint16(header[28:30])
		wantFraction := uint16(wantTime.Nanosecond() / int(waveform.MiniSeedTimeCorrectionUnit))
		if gotFraction != wantFraction {
			t.Fatalf("record %d BTIME fraction = %d, want %d (100us units)", i, gotFraction, wantFraction)
		}
		wantCorrection := int32(-jitters[i] / waveform.MiniSeedTimeCorrectionUnit)
		if series.FixedSection.TimeCorrection != wantCorrection {
			t.Fatalf("record %d time correction = %d, want %d", i, series.FixedSection.TimeCorrection, wantCorrection)
		}
		if series.FixedSection.ActivityFlags&mseedio.ACTIVITY_TIME_CORRECTION_APPLIED == 0 {
			t.Fatalf("record %d does not mark time correction as applied", i)
		}
		reconstructedRecordTime := wantTime.Add(-time.Duration(series.FixedSection.TimeCorrection) * waveform.MiniSeedTimeCorrectionUnit)
		if !reconstructedRecordTime.Equal(recordTime) {
			t.Fatalf("record %d reconstructed raw time = %s, want %s", i, reconstructedRecordTime, recordTime)
		}
		if len(series.BlocketteSections) != 2 || series.BlocketteSections[1].BlocketteCode != 1001 || series.BlocketteSections[1].TimingQuality != 90 {
			t.Fatalf("record %d blockettes = %#v", i, series.BlocketteSections)
		}
	}
}
