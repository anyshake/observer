// Package groundmotion derives peak ground acceleration, velocity and displacement
// from raw counts of a velocity or acceleration channel.
package groundmotion

import (
	"math"
	"time"
)

const (
	UNIT_VELOCITY     = "m/s"
	UNIT_ACCELERATION = "m/s2"

	// HIGHPASS_CORNER removes DC offset and drift before integration.
	HIGHPASS_CORNER = 0.1
	// GAP_TOLERANCE resets the filters when samples are not contiguous.
	GAP_TOLERANCE = 2 * time.Second
	// SETTLE_DURATION hides the filter start-up transient after a reset.
	SETTLE_DURATION = 10 * time.Second
)

// Peak holds absolute peak values in SI units: m/s², m/s and m.
type Peak struct {
	PGA float64
	PGV float64
	PGD float64
}

// highpass is a first-order causal high-pass filter.
type highpass struct {
	alpha            float64
	previous, output float64
	primed           bool
}

func newHighpass(corner, dt float64) highpass {
	return highpass{alpha: 1 / (1 + 2*math.Pi*corner*dt)}
}

func (h *highpass) apply(value float64) float64 {
	if !h.primed {
		h.previous, h.output, h.primed = value, 0, true
		return 0
	}
	h.output = h.alpha * (h.output + value - h.previous)
	h.previous = value
	return h.output
}

// integrator integrates by the trapezoidal rule and high-pass filters the result
// so that integration drift does not accumulate.
type integrator struct {
	filter   highpass
	previous float64
	sum      float64
}

func (i *integrator) apply(value, dt float64) float64 {
	i.sum += (value + i.previous) * dt / 2
	i.previous = value
	return i.filter.apply(i.sum)
}

// Channel keeps the filter state of a single channel between consecutive blocks.
type Channel struct {
	gain float64
	unit string

	poles     []complex128
	frequency float64
	corrector *correction

	rate    int
	end     time.Time
	settled time.Time

	input        highpass
	first        integrator
	second       integrator
	lastVelocity float64
	hasVelocity  bool
}

// NewChannel creates a calculator for a channel with the given sensitivity
// (counts per SI unit). Unsupported units or gains return nil.
func NewChannel(gain float64, unit string) *Channel {
	if gain <= 0 || math.IsNaN(gain) || math.IsInf(gain, 0) {
		return nil
	}
	if unit != UNIT_VELOCITY && unit != UNIT_ACCELERATION {
		return nil
	}
	return &Channel{gain: gain, unit: unit}
}

// WithResponse enables instrument correction for a channel whose sensor has a
// low-frequency roll-off (e.g. a geophone). poles are the analog sensor poles
// in rad/s and frequency is the frequency (Hz) at which gain is defined.
// Poles that need no correction are ignored, so flat responses are unaffected.
func (c *Channel) WithResponse(poles []complex128, frequency float64) *Channel {
	if c == nil {
		return nil
	}
	c.poles = lowFrequencyPoles(poles)
	c.frequency = frequency
	return c
}

func (c *Channel) reset(start time.Time, rate int) {
	dt := 1 / float64(rate)
	c.rate = rate
	c.settled = start.Add(SETTLE_DURATION)
	c.input = newHighpass(HIGHPASS_CORNER, dt)
	c.first = integrator{filter: newHighpass(HIGHPASS_CORNER, dt)}
	c.second = integrator{filter: newHighpass(HIGHPASS_CORNER, dt)}
	c.lastVelocity, c.hasVelocity = 0, false
	c.corrector = nil
	if len(c.poles) > 0 {
		c.corrector = newCorrection(c.poles, c.frequency, rate)
	}
}

// Process consumes a block of raw counts starting at start and returns the
// peaks of the block. ok is false while the filters settle after a reset.
func (c *Channel) Process(start time.Time, rate int, samples []int32) (peak Peak, ok bool) {
	if c == nil || rate < 1 || len(samples) == 0 {
		return Peak{}, false
	}
	if c.rate != rate || c.end.IsZero() || start.Sub(c.end).Abs() > GAP_TOLERANCE {
		c.reset(start, rate)
	}
	dt := 1 / float64(rate)
	for _, sample := range samples {
		value := c.input.apply(float64(sample) / c.gain)
		if c.corrector != nil {
			value = c.corrector.apply(value)
		}
		var acceleration, velocity, displacement float64
		if c.unit == UNIT_ACCELERATION {
			acceleration = value
			velocity = c.first.apply(acceleration, dt)
			displacement = c.second.apply(velocity, dt)
		} else {
			velocity = value
			if c.hasVelocity {
				acceleration = (velocity - c.lastVelocity) / dt
			}
			c.lastVelocity, c.hasVelocity = velocity, true
			displacement = c.first.apply(velocity, dt)
		}
		peak.PGA = max(peak.PGA, math.Abs(acceleration))
		peak.PGV = max(peak.PGV, math.Abs(velocity))
		peak.PGD = max(peak.PGD, math.Abs(displacement))
	}
	c.end = start.Add(time.Duration(len(samples)) * time.Second / time.Duration(rate))
	return peak, !c.end.Before(c.settled)
}
