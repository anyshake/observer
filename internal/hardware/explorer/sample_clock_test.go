package explorer

import (
	"math/rand"
	"testing"
	"time"
)

// runSampleClock feeds one-second packets of 250 samples whose measured start
// is trueTime(i) plus receive noise, and checks the continuity invariant that
// SeedLink relies on: each packet starts where the previous one ended, to
// within the bounded slew and far below one sample period.
func runSampleClock(t *testing.T, packets int, trueTime func(int) time.Time, noise func(int) time.Duration) (*sampleClock, []time.Time) {
	t.Helper()
	const rate, samples = 250, 250
	var clock sampleClock
	published := make([]time.Time, packets)
	maxStep := time.Second * NTP_TIMELINE_MAX_SLEW_PPM / 1_000_000
	for i := range packets {
		got, event, _ := clock.Next(trueTime(i).Add(noise(i)), rate, samples, true)
		if i > 0 && event != sampleClockContinued {
			t.Fatalf("packet %d restarted the timeline (event %d)", i, event)
		}
		if i > 0 {
			step := got.Sub(published[i-1]) - time.Second
			if step < -maxStep-time.Nanosecond || step > maxStep+time.Nanosecond {
				t.Fatalf("packet %d moved the timeline by %s, more than %s", i, step, maxStep)
			}
		}
		published[i] = got
	}
	return &clock, published
}

func TestSampleClockAbsorbsReceiveJitter(t *testing.T) {
	start := time.Unix(1_800_000_000, 0)
	rng := rand.New(rand.NewSource(1))
	_, published := runSampleClock(t, 3600,
		func(i int) time.Time { return start.Add(time.Duration(i) * time.Second) },
		func(int) time.Duration { return time.Duration(rng.Intn(16_000_000)) - 8*time.Millisecond },
	)
	// The first packet anchors on a noisy read; at 100 ppm an 8 ms anchor
	// error is absorbed within a few minutes.
	for i := 300; i < len(published); i++ {
		if d := published[i].Sub(start.Add(time.Duration(i) * time.Second)); d < -3*time.Millisecond || d > 3*time.Millisecond {
			t.Fatalf("packet %d is %s from true time", i, d)
		}
	}
}

func TestSampleClockFollowsCrystalDrift(t *testing.T) {
	// The ADC clock runs 50 ppm fast: 250 samples span 999.95 ms of true time.
	start := time.Unix(1_800_000_000, 0)
	rng := rand.New(rand.NewSource(2))
	period := time.Second - 50*time.Microsecond
	_, published := runSampleClock(t, 86400,
		func(i int) time.Time { return start.Add(time.Duration(i) * period) },
		func(int) time.Duration { return time.Duration(rng.Intn(10_000_000)) - 5*time.Millisecond },
	)
	for i := 3600; i < len(published); i++ {
		if d := published[i].Sub(start.Add(time.Duration(i) * period)); d < -3*time.Millisecond || d > 3*time.Millisecond {
			t.Fatalf("packet %d drifted %s from true time", i, d)
		}
	}
}

func TestSampleClockSlewsNTPStep(t *testing.T) {
	// The NTP-disciplined time source steps 20 ms at packet 600.
	start := time.Unix(1_800_000_000, 0)
	measured := func(i int) time.Time {
		at := start.Add(time.Duration(i) * time.Second)
		if i >= 600 {
			at = at.Add(20 * time.Millisecond)
		}
		return at
	}
	_, published := runSampleClock(t, 1200, measured, func(int) time.Duration { return 0 })
	if d := published[1199].Sub(measured(1199)); d < -time.Millisecond || d > time.Millisecond {
		t.Fatalf("timeline is %s from the stepped time after 10 minutes", d)
	}
}

func TestSampleClockResetsOnPersistentLargeStep(t *testing.T) {
	start := time.Unix(1_800_000_000, 0)
	var clock sampleClock
	clock.Next(start, 250, 250, true)
	// One very late read is not a step.
	if got, event, _ := clock.Next(start.Add(time.Second+300*time.Millisecond), 250, 250, true); event != sampleClockContinued || !got.Equal(start.Add(time.Second)) {
		t.Fatalf("single late read: got %s, event %d", got, event)
	}
	if _, event, _ := clock.Next(start.Add(2*time.Second), 250, 250, true); event != sampleClockContinued {
		t.Fatalf("recovered read: event %d", event)
	}
	// A persistent 1 s step restarts the timeline at the measured time.
	var event sampleClockEvent
	var got, at time.Time
	for i := 3; event == sampleClockContinued; i++ {
		if i > 10 {
			t.Fatal("persistent step was never applied")
		}
		at = start.Add(time.Duration(i+1) * time.Second)
		got, event, _ = clock.Next(at, 250, 250, true)
	}
	if event != sampleClockStepped || !got.Equal(at) {
		t.Fatalf("persistent step: got %s, event %d, want %s", got, event, at)
	}
}

func TestSampleClockRestartsOnKnownGap(t *testing.T) {
	start := time.Unix(1_800_000_000, 0)
	var clock sampleClock
	clock.Next(start, 250, 25, true)
	at := start.Add(time.Second)
	if got, event, _ := clock.Next(at, 250, 25, false); event != sampleClockStarted || !got.Equal(at) {
		t.Fatalf("gap: got %s, event %d", got, event)
	}
	if got, event, _ := clock.Next(at.Add(time.Second), 100, 10, true); event != sampleClockStarted || !got.Equal(at.Add(time.Second)) {
		t.Fatalf("rate change: got %s, event %d", got, event)
	}
}
