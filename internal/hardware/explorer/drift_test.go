package explorer

import (
	"math"
	"testing"
	"time"

	"github.com/anyshake/observer/pkg/ringbuf"
)

func TestLongTermClockDriftPPM(t *testing.T) {
	t.Parallel()

	base := time.Unix(1700000000, 0)
	tests := []struct {
		name    string
		values  []clockDrift
		wantPPM float64
		wantLen int
	}{
		{name: "empty"},
		{name: "single measurement", values: []clockDrift{{measuredAt: base}}, wantLen: 1},
		{name: "positive drift", values: []clockDrift{{measuredAt: base}, {measuredAt: base.Add(time.Hour), offset: 36 * time.Millisecond}}, wantPPM: 10, wantLen: 2},
		{name: "negative drift", values: []clockDrift{{measuredAt: base}, {measuredAt: base.Add(time.Hour), offset: -36 * time.Millisecond}}, wantPPM: -10, wantLen: 2},
		{name: "equal timestamps", values: []clockDrift{{measuredAt: base}, {measuredAt: base, offset: time.Millisecond}}, wantLen: 2},
		{name: "backward timestamps", values: []clockDrift{{measuredAt: base}, {measuredAt: base.Add(-time.Minute)}}, wantLen: 2},
		{name: "insufficient recent data", values: []clockDrift{{measuredAt: base}, {measuredAt: base.Add(2 * time.Hour)}}, wantLen: 2},
		{
			name: "prune old data and include cutoff",
			values: []clockDrift{
				{measuredAt: base, offset: time.Second},
				{measuredAt: base.Add(time.Hour)},
				{measuredAt: base.Add(2 * time.Hour), offset: 36 * time.Millisecond},
			},
			wantPPM: 10, wantLen: 2,
		},
		{name: "positive outlier resets history", values: []clockDrift{{measuredAt: base}, {measuredAt: base.Add(time.Second), offset: 50 * time.Microsecond}}},
		{name: "negative outlier resets history", values: []clockDrift{{measuredAt: base}, {measuredAt: base.Add(time.Second), offset: -50 * time.Microsecond}}},
		{name: "below outlier threshold", values: []clockDrift{{measuredAt: base}, {measuredAt: base.Add(time.Second), offset: 49 * time.Microsecond}}, wantPPM: 49, wantLen: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := ringbuf.New[clockDrift](10)
			buf.Push(tt.values...)
			if got := getLongTermClockDriftPPM(buf, time.Hour); math.IsNaN(got) || math.Abs(got-tt.wantPPM) > 1e-9 {
				t.Errorf("clock drift = %v PPM, want %v", got, tt.wantPPM)
			}
			if got := buf.Len(); got != tt.wantLen {
				t.Errorf("history length = %d, want %d", got, tt.wantLen)
			}
		})
	}
}
