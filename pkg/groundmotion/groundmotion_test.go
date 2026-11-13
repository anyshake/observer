package groundmotion

import (
	"math"
	"testing"
	"time"
)

func sine(gain, amplitude, frequency float64, rate, offset, n int) []int32 {
	samples := make([]int32, n)
	for i := range samples {
		t := float64(offset+i) / float64(rate)
		samples[i] = int32(math.Round(gain * amplitude * math.Sin(2*math.Pi*frequency*t)))
	}
	return samples
}

func run(t *testing.T, c *Channel, gain, amplitude, frequency float64) Peak {
	t.Helper()
	const rate = 200
	start := time.Unix(1000, 0)
	var last Peak
	for second := 0; second < 30; second++ {
		peak, ok := c.Process(start.Add(time.Duration(second)*time.Second), rate, sine(gain, amplitude, frequency, rate, second*rate, rate))
		if second < 9 && ok {
			t.Fatal("peak reported before filters settled")
		}
		if ok {
			last = peak
		}
	}
	return last
}

func near(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > want*0.1 {
		t.Errorf("%s = %g, want about %g", name, got, want)
	}
}

func TestAccelerationChannel(t *testing.T) {
	const gain, amplitude, frequency = 26104.7, 0.5, 2.0
	peak := run(t, NewChannel(gain, UNIT_ACCELERATION), gain, amplitude, frequency)
	omega := 2 * math.Pi * frequency
	near(t, "PGA", peak.PGA, amplitude)
	near(t, "PGV", peak.PGV, amplitude/omega)
	near(t, "PGD", peak.PGD, amplitude/omega/omega)
}

func TestVelocityChannel(t *testing.T) {
	const gain, amplitude, frequency = 85899345880, 0.001, 2.0
	peak := run(t, NewChannel(gain, UNIT_VELOCITY), gain, amplitude, frequency)
	omega := 2 * math.Pi * frequency
	near(t, "PGA", peak.PGA, amplitude*omega)
	near(t, "PGV", peak.PGV, amplitude)
	near(t, "PGD", peak.PGD, amplitude/omega)
}

func TestGapResetsAndInvalidChannels(t *testing.T) {
	c := NewChannel(1000, UNIT_ACCELERATION)
	start := time.Unix(1000, 0)
	for second := 0; second < 12; second++ {
		c.Process(start.Add(time.Duration(second)*time.Second), 100, make([]int32, 100))
	}
	if _, ok := c.Process(start.Add(time.Minute), 100, make([]int32, 100)); ok {
		t.Fatal("gap did not reset settling")
	}
	if NewChannel(0, UNIT_VELOCITY) != nil || NewChannel(1, "counts") != nil {
		t.Fatal("invalid channel accepted")
	}
	var nilChannel *Channel
	if _, ok := nilChannel.Process(start, 100, []int32{1}); ok {
		t.Fatal("nil channel produced a peak")
	}
}
