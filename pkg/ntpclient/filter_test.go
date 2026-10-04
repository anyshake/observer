package ntpclient

import (
	"testing"
	"time"
)

func TestCombineSamples(t *testing.T) {
	t.Parallel()
	largeOffset := 30 * 365 * 24 * time.Hour
	for _, tt := range []struct {
		name       string
		samples    []clockSample
		target     int
		wantOffset time.Duration
		wantServer string
		wantErr    bool
	}{
		{
			name:    "single source",
			samples: []clockSample{{"only", 1234567 * time.Nanosecond, time.Millisecond}},
			target:  5, wantOffset: 1234567 * time.Nanosecond, wantServer: "only",
		},
		{
			name: "nanosecond precision with decades-long offset",
			samples: []clockSample{
				{"first", largeOffset + 123*time.Nanosecond, time.Millisecond},
				{"second", largeOffset + 127*time.Nanosecond, time.Millisecond},
			},
			target: 5, wantOffset: largeOffset + 125*time.Nanosecond, wantServer: "first",
		},
		{
			name: "negative offsets retain precision",
			samples: []clockSample{
				{"first", -largeOffset - 123*time.Nanosecond, time.Millisecond},
				{"second", -largeOffset - 127*time.Nanosecond, time.Millisecond},
			},
			target: 5, wantOffset: -largeOffset - 125*time.Nanosecond, wantServer: "first",
		},
		{
			name: "reject low-latency outlier before limiting sample count",
			samples: []clockSample{
				{"wrong", time.Second, time.Microsecond},
				{"good", 10 * time.Millisecond, time.Millisecond},
				{"also-good", 12 * time.Millisecond, 2 * time.Millisecond},
			},
			target: 1, wantOffset: 10 * time.Millisecond, wantServer: "good",
		},
		{
			name: "weight by synchronization distance",
			samples: []clockSample{
				{"noisy", 12 * time.Millisecond, 2 * time.Millisecond},
				{"precise", 10 * time.Millisecond, time.Millisecond},
			},
			target: 5, wantOffset: 10400 * time.Microsecond, wantServer: "precise",
		},
		{
			name: "clamp estimate to common interval",
			samples: []clockSample{
				{"first", 0, time.Millisecond},
				{"second", 3 * time.Millisecond, 2 * time.Millisecond},
			},
			target: 5, wantOffset: time.Millisecond, wantServer: "first",
		},
		{
			name: "sample limit preserves majority intersection",
			samples: []clockSample{
				{"first", 0, time.Millisecond},
				{"second", 3 * time.Millisecond, 2 * time.Millisecond},
			},
			target: 1, wantOffset: time.Millisecond, wantServer: "first",
		},
		{
			name: "split vote has no majority",
			samples: []clockSample{
				{"first", 0, time.Millisecond},
				{"second", 0, time.Millisecond},
				{"third", time.Second, time.Millisecond},
				{"fourth", time.Second, time.Millisecond},
			},
			target: 5, wantErr: true,
		},
		{
			name: "disagreeing pair",
			samples: []clockSample{
				{"first", 0, time.Millisecond},
				{"second", time.Second, time.Millisecond},
			},
			target: 5, wantErr: true,
		},
		{name: "no samples", target: 5, wantErr: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			offset, server, err := combineSamples(tt.samples, tt.target)
			if (err != nil) != tt.wantErr {
				t.Fatalf("combineSamples() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (offset != tt.wantOffset || server != tt.wantServer) {
				t.Errorf("combineSamples() = (%v, %q), want (%v, %q)", offset, server, tt.wantOffset, tt.wantServer)
			}
		})
	}
}
