package explorer

import (
	"math"
	"time"
)

// sampleClock publishes NTP-mode packet times on a continuous sample-count
// timeline. Per-packet measurements are noisy (serial latency) and step when
// the NTP-disciplined time source is updated; publishing them directly makes
// adjacent records overlap or leave gaps, and downstream SeedLink clients drop
// the overlapping samples. Instead, the timeline advances exactly by the number
// of samples and is slewed toward the filtered measurement at a bounded rate.
// Only a persistent large error (a real time step or lost data) restarts it.
type sampleClock struct {
	initialized   bool
	sampleRate    int
	anchor        time.Time
	samples       int64
	slew          time.Duration
	filteredError float64       // nanoseconds, measured minus published
	outlierTime   time.Duration // Consecutive data duration beyond the reset threshold.
}

type sampleClockEvent int

const (
	sampleClockContinued sampleClockEvent = iota
	sampleClockStarted                    // First packet, rate change or caller-detected discontinuity.
	sampleClockStepped                    // Persistent error beyond NTP_TIMELINE_RESET_THRESHOLD.
)

func (c *sampleClock) Reset() {
	*c = sampleClock{}
}

// Next returns the publish time of a packet that starts at the measured time
// and holds the given number of samples per channel. Set contiguous to false
// when the caller knows packets are missing; the timeline then restarts at the
// measurement. The returned error is measured minus the continuous prediction.
func (c *sampleClock) Next(measured time.Time, sampleRate, samples int, contiguous bool) (time.Time, sampleClockEvent, time.Duration) {
	if sampleRate <= 0 || samples <= 0 {
		c.Reset()
		return measured, sampleClockStarted, 0
	}
	if !c.initialized || c.sampleRate != sampleRate || !contiguous {
		c.start(measured, sampleRate, samples)
		return measured, sampleClockStarted, 0
	}

	predicted := c.anchor.Add(samplesDuration(c.samples, sampleRate) + c.slew)
	measuredError := measured.Sub(predicted)
	if measuredError > NTP_TIMELINE_RESET_THRESHOLD || measuredError < -NTP_TIMELINE_RESET_THRESHOLD {
		// A burst of late reads must not restart the timeline. Lost data or a
		// real clock step keeps the error beyond the threshold.
		c.outlierTime += samplesDuration(int64(samples), sampleRate)
		if c.outlierTime >= NTP_TIMELINE_RESET_AFTER {
			c.start(measured, sampleRate, samples)
			return measured, sampleClockStepped, measuredError
		}
	} else {
		c.outlierTime = 0
		weight := math.Min(1, float64(samplesDuration(int64(samples), sampleRate))/float64(NTP_TIMELINE_FILTER_WINDOW))
		c.filteredError += weight * (float64(measuredError) - c.filteredError)
	}
	c.samples += int64(samples)

	// Bound the change between adjacent packets to NTP_TIMELINE_MAX_SLEW_PPM of
	// their duration, far below one sample period at any supported rate.
	maxStep := float64(samplesDuration(int64(samples), sampleRate)) * NTP_TIMELINE_MAX_SLEW_PPM / 1e6
	step := math.Max(-maxStep, math.Min(maxStep, c.filteredError))
	c.slew += time.Duration(math.Round(step))
	c.filteredError -= step
	return predicted, sampleClockContinued, measuredError
}

func (c *sampleClock) start(measured time.Time, sampleRate, samples int) {
	*c = sampleClock{
		initialized: true,
		sampleRate:  sampleRate,
		anchor:      measured,
		samples:     int64(samples),
	}
}

func samplesDuration(samples int64, sampleRate int) time.Duration {
	rate := int64(sampleRate)
	return time.Duration(samples/rate)*time.Second + time.Duration(samples%rate)*time.Second/time.Duration(rate)
}
