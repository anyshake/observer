package waveform

import (
	"fmt"
	"time"
)

// MiniSeedTimeCorrectionUnit is the fixed-header correction resolution.
const MiniSeedTimeCorrectionUnit = 100 * time.Microsecond

// Timeline places continuous packets on the sample timeline while using
// adjacent raw timestamps to distinguish clock drift from a real discontinuity.
type Timeline struct {
	startTime      time.Time
	lastRecordTime time.Time
	sampleRate     int
	lastSamples    int
	samples        int64
}

func (t *Timeline) AppendStartTime(recordTime time.Time, sampleRate, samples int, allowedJitter time.Duration) time.Time {
	recordTime = recordTime.UTC()
	startTime := recordTime
	continuous := false
	if t.sampleRate == sampleRate && sampleRate > 0 {
		expectedRecordTime := t.lastRecordTime.Add(sampleOffset(int64(t.lastSamples), sampleRate))
		jitter := recordTime.Sub(expectedRecordTime)
		if jitter > -allowedJitter && jitter < allowedJitter {
			startTime = t.startTime.Add(sampleOffset(t.samples, sampleRate))
			continuous = true
		}
	}
	if !continuous {
		t.startTime = startTime
		t.sampleRate = sampleRate
		t.samples = 0
	}
	t.lastRecordTime = recordTime
	t.lastSamples = samples
	t.samples += int64(samples)
	return startTime
}

func sampleOffset(samples int64, sampleRate int) time.Duration {
	rate := int64(sampleRate)
	seconds := samples / rate
	remainder := samples % rate
	return time.Duration(seconds)*time.Second + time.Duration(remainder)*time.Second/time.Duration(rate)
}

// MiniSeedTimeCorrection converts the difference between a raw record time and
// its effective start time to the signed MiniSEED 2 fixed-header field.
func MiniSeedTimeCorrection(recordTime, correctedStartTime time.Time) (int32, error) {
	correction := correctedStartTime.Sub(recordTime.UTC())
	if correction%MiniSeedTimeCorrectionUnit != 0 {
		return 0, fmt.Errorf("MiniSEED time correction %s is not representable in 0.0001 second units", correction)
	}
	units := correction / MiniSeedTimeCorrectionUnit
	if units < -1<<31 || units > 1<<31-1 {
		return 0, fmt.Errorf("MiniSEED time correction %s exceeds the int32 field range", correction)
	}
	return int32(units), nil
}
