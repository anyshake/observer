package waveform

import "time"

const (
	// TimelineMaxSlewPPM bounds how fast a continuous segment follows its
	// recorded times. It matches the source-side NTP slew, so a recorded time
	// that is already continuous is reproduced within millisecond rounding.
	TimelineMaxSlewPPM = 100
	// TimelineTolerance absorbs millisecond rounding and the receive jitter of
	// records stored before source-side disciplining. The source itself only
	// steps by more than its 50 ms reset threshold, which stays visible here.
	TimelineTolerance = 20 * time.Millisecond
)

// Timeline places records on a continuous sample timeline. Recorded times are
// already disciplined at the source, so this only removes millisecond rounding
// and detects discontinuities; it never accumulates an offset from them.
type Timeline struct {
	startTime  time.Time
	sampleRate int
	samples    int
}

// AppendStartTime returns the start time to write for a record. A record that
// continues the previous one within tolerance keeps the sample timeline and
// moves toward its recorded time by at most TimelineMaxSlewPPM; otherwise a new
// segment starts at the recorded time.
func (t *Timeline) AppendStartTime(recordTime time.Time, sampleRate, samples int) time.Time {
	recordTime = recordTime.UTC()
	startTime := recordTime
	if t.sampleRate == sampleRate && sampleRate > 0 && !t.startTime.IsZero() {
		expected := t.startTime.Add(SampleOffset(int64(t.samples), sampleRate))
		offset := recordTime.Sub(expected)
		if offset > -TimelineTolerance && offset < TimelineTolerance {
			maxStep := SampleOffset(int64(samples), sampleRate) * TimelineMaxSlewPPM / 1_000_000
			startTime = expected.Add(min(maxStep, max(-maxStep, offset)))
		}
	}
	t.startTime = startTime
	t.sampleRate = sampleRate
	t.samples = samples
	return startTime
}

// SampleOffset returns the time spanned by the given number of samples.
func SampleOffset(samples int64, sampleRate int) time.Duration {
	rate := int64(sampleRate)
	seconds := samples / rate
	remainder := samples % rate
	return time.Duration(seconds)*time.Second + time.Duration(remainder)*time.Second/time.Duration(rate)
}
