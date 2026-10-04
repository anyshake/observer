package timesource_test

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/anyshake/observer/pkg/timesource"
)

func TestSourceUsesClockAndReturnsUTC(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.FixedZone("UTC+8", 8*60*60))
	source := timesource.New(func() time.Time { return now })
	now = now.Add(1500 * time.Millisecond)
	if got := source.Now(); !got.Equal(now) || got.Location() != time.UTC {
		t.Fatalf("Now() = %v, want %v in UTC", got, now.UTC())
	}
}

func TestSourceUpdateAppliesOffsetAndDrift(t *testing.T) {
	t.Parallel()
	for _, ppm := range []float64{-10, 0, 10} {
		local := time.Unix(1000, 0)
		now := local
		source := timesource.New(func() time.Time { return now })
		reference := local.Add(5 * time.Hour)
		source.Update(local, reference, ppm, nil)
		if got := source.Now(); !got.Equal(reference) {
			t.Fatalf("Now() immediately after Update() = %v, want %v", got, reference)
		}
		now = local.Add(100 * time.Second)
		want := reference.Add(100*time.Second + time.Duration(ppm)*100*time.Microsecond)
		if got := source.Now(); !got.Equal(want) {
			t.Errorf("drift %v PPM: Now() = %v, want %v", ppm, got, want)
		}
	}
}

func TestSourceUpdateCanReplaceClock(t *testing.T) {
	t.Parallel()
	oldClock := time.Unix(10, 0)
	newClock := time.Unix(100, 0)
	reference := time.Unix(10000, 0)
	source := timesource.New(func() time.Time { return oldClock })
	source.Update(newClock, reference, 0, func() time.Time { return newClock })
	oldClock = oldClock.Add(time.Hour)
	newClock = newClock.Add(2 * time.Second)
	if got := source.Now(); !got.Equal(reference.Add(2 * time.Second)) {
		t.Fatalf("replacement clock ignored: %v", got)
	}
}

func TestSourceDefaultsToSystemTime(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		source := timesource.New(nil)
		time.Sleep(time.Second)
		if got := source.Now(); !got.Equal(time.Now()) {
			t.Fatalf("default Now() = %v, want %v", got, time.Now())
		}
	})
}
