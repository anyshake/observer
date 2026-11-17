package waveform

import (
	"testing"
	"time"
)

func TestTimelineRemovesMillisecondRounding(t *testing.T) {
	// 30 Hz records of 30 samples: the true start of each record is exact, but
	// rounding to milliseconds moves it by up to 0.5 ms.
	start := time.Date(2026, 11, 6, 3, 56, 40, 123000000, time.UTC)
	var timeline Timeline
	for i := 0; i < 10; i++ {
		exact := start.Add(time.Duration(i) * 1000 * time.Millisecond).Add(time.Duration(i%3) * 333333)
		got := timeline.AppendStartTime(exact.Truncate(time.Millisecond), 30, 30)
		want := start.Add(time.Duration(i) * time.Second)
		// The timeline may follow rounding by at most 100 ppm per record.
		if d := got.Sub(want); d < 0 || d > time.Duration(i)*100*time.Microsecond {
			t.Fatalf("record %d start = %s, want continuous %s", i, got, want)
		}
	}
}

func TestTimelineDoesNotAccumulateDrift(t *testing.T) {
	// Source-disciplined times run 50 ppm slower than the nominal rate. The
	// timeline must follow them instead of drifting away for a whole day.
	start := time.Date(2026, 11, 6, 0, 0, 0, 0, time.UTC)
	var timeline Timeline
	for i := 0; i < 86400; i++ {
		recorded := start.Add(time.Duration(i) * (time.Second + 50*time.Microsecond)).Truncate(time.Millisecond)
		got := timeline.AppendStartTime(recorded, 100, 100)
		if d := got.Sub(recorded); d <= -2*time.Millisecond || d >= 2*time.Millisecond {
			t.Fatalf("record %d start = %s, recorded %s", i, got, recorded)
		}
	}
}

func TestTimelinePreservesDiscontinuities(t *testing.T) {
	start := time.Date(2026, 11, 6, 3, 56, 40, 123000000, time.UTC)
	var timeline Timeline
	timeline.AppendStartTime(start, 100, 100)
	for _, tc := range []struct {
		name string
		at   time.Time
		rate int
	}{
		{"time step", start.Add(time.Second + 60*time.Millisecond), 100},
		{"gap", start.Add(3 * time.Second), 100},
		{"overlap", start.Add(3*time.Second + 900*time.Millisecond), 100},
		{"rate change", start.Add(4 * time.Second), 50},
	} {
		if got := timeline.AppendStartTime(tc.at, tc.rate, tc.rate); !got.Equal(tc.at) {
			t.Fatalf("%s was hidden: got %s, want %s", tc.name, got, tc.at)
		}
	}
}
