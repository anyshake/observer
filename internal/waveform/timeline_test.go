package waveform

import (
	"testing"
	"time"
)

func TestTimelineCorrectsDriftAndPreservesTimeTear(t *testing.T) {
	start := time.Date(2026, 11, 6, 3, 56, 40, 123000000, time.UTC)
	var timeline Timeline

	if got := timeline.AppendStartTime(start, 100, 100, 20*time.Millisecond); !got.Equal(start) {
		t.Fatalf("first start = %s, want %s", got, start)
	}
	drifted := start.Add(time.Second + 15*time.Millisecond)
	corrected := timeline.AppendStartTime(drifted, 100, 100, 20*time.Millisecond)
	if want := start.Add(time.Second); !corrected.Equal(want) {
		t.Fatalf("corrected start = %s, want %s", corrected, want)
	}
	correction, err := MiniSeedTimeCorrection(drifted, corrected)
	if err != nil || correction != -150 {
		t.Fatalf("correction = %d, %v; want -150", correction, err)
	}

	timeTear := start.Add(3 * time.Second)
	if got := timeline.AppendStartTime(timeTear, 100, 100, 20*time.Millisecond); !got.Equal(timeTear) {
		t.Fatalf("time tear was hidden: got %s, want %s", got, timeTear)
	}
}
