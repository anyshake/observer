package groundmotion

import (
	"math"
	"math/cmplx"
	"testing"
	"time"
)

// E-C121G geophone (LGT-4.5C) poles in rad/s; sensitivity defined at 4.5 Hz.
var geophonePoles = []complex128{
	complex(-2.2211060060879837, 2.221776881419294),
	complex(-2.2211060060879837, -2.221776881419294),
}

const geophoneGain, geophoneFrequency = 85899345880.0, 4.5

// geophoneResponse is T(f)/|T(f0)| for T(s) = s² / Π(s - p).
func geophoneResponse(frequency float64) complex128 {
	eval := func(f float64) complex128 {
		s := complex(0, 2*math.Pi*f)
		value := s * s
		for _, pole := range geophonePoles {
			value /= s - pole
		}
		return value
	}
	return eval(frequency) / complex(cmplx.Abs(eval(geophoneFrequency)), 0)
}

// geophoneCounts is the steady-state output for a velocity sine of amplitude
// (m/s) and frequency (Hz).
func geophoneCounts(amplitude, frequency float64, rate, offset, n int) []int32 {
	response := geophoneResponse(frequency)
	samples := make([]int32, n)
	for i := range samples {
		t := float64(offset+i) / float64(rate)
		value := geophoneGain * amplitude * cmplx.Abs(response) * math.Sin(2*math.Pi*frequency*t+cmplx.Phase(response))
		samples[i] = int32(math.Round(value))
	}
	return samples
}

func peakVelocity(c *Channel, amplitude, frequency float64) float64 {
	const rate = 100
	start := time.Unix(1000, 0)
	var last float64
	for second := range 90 {
		peak, ok := c.Process(start.Add(time.Duration(second)*time.Second), rate, geophoneCounts(amplitude, frequency, rate, second*rate, rate))
		if ok && second > 60 {
			last = max(last, peak.PGV)
		}
	}
	return last
}

func TestResponseCorrectionRestoresGeophoneRollOff(t *testing.T) {
	const amplitude, frequency = 1e-4, 0.3 // well below the 0.5 Hz natural frequency
	raw := peakVelocity(NewChannel(geophoneGain, UNIT_VELOCITY), amplitude, frequency)
	corrected := peakVelocity(NewChannel(geophoneGain, UNIT_VELOCITY).WithResponse(geophonePoles, geophoneFrequency), amplitude, frequency)

	// The only remaining attenuation is the 0.1 Hz first-order input high-pass.
	want := amplitude * 0.3 / math.Hypot(0.3, HIGHPASS_CORNER)
	if math.Abs(corrected-want) > want*0.05 {
		t.Fatalf("corrected PGV = %g, want about %g", corrected, want)
	}
	if raw > want*0.5 {
		t.Fatalf("uncorrected PGV = %g is not attenuated; test does not exercise the roll-off", raw)
	}
}

func TestResponseCorrectionKeepsPassBand(t *testing.T) {
	const amplitude, frequency = 1e-4, 4.5
	corrected := peakVelocity(NewChannel(geophoneGain, UNIT_VELOCITY).WithResponse(geophonePoles, geophoneFrequency), amplitude, frequency)
	if math.Abs(corrected-amplitude) > amplitude*0.03 {
		t.Fatalf("PGV at the calibration frequency = %g, want %g", corrected, amplitude)
	}
}

func TestResponseCorrectionIgnoresFlatAndHighPoles(t *testing.T) {
	// A 27 Hz anti-alias pole is not inverted, and flat sensors get no filter.
	if c := NewChannel(1, UNIT_VELOCITY).WithResponse([]complex128{-169.646}, 1); len(c.poles) != 0 {
		t.Fatalf("high-frequency pole selected for correction: %v", c.poles)
	}
	if newCorrection(nil, 1, 100) != nil {
		t.Fatal("correction designed without poles")
	}
	if newCorrection([]complex128{complex(0.5, 0)}, 1, 100) != nil {
		t.Fatal("unstable pole accepted")
	}
	var nilChannel *Channel
	if nilChannel.WithResponse(geophonePoles, 1) != nil {
		t.Fatal("nil channel became non-nil")
	}
}

func TestResponseCorrectionStability(t *testing.T) {
	c := newCorrection(geophonePoles, geophoneFrequency, 100)
	if c == nil {
		t.Fatal("no correction designed")
	}
	// The impulse response must decay: all replacement poles lie inside the unit circle.
	value := c.apply(1)
	for range 100 * 120 {
		value = c.apply(0)
	}
	if math.IsNaN(value) || math.Abs(value) > 1e-6 {
		t.Fatalf("impulse response did not decay: %g", value)
	}
}

func BenchmarkProcessCorrectedChannel(b *testing.B) {
	c := NewChannel(geophoneGain, UNIT_VELOCITY).WithResponse(geophonePoles, geophoneFrequency)
	samples := geophoneCounts(1e-4, 1, 100, 0, 100)
	start := time.Unix(1000, 0)
	b.ReportAllocs()
	for i := 0; b.Loop(); i++ {
		c.Process(start.Add(time.Duration(i)*time.Second), 100, samples)
	}
}
