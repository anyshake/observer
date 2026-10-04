package ntpclient

import (
	"cmp"
	"errors"
	"math"
	"slices"
	"time"
)

// Each sample defines an offset interval using its synchronization distance.
// Require a strict majority to agree before choosing the lowest-distance
// measurements. A single usable source is accepted without cross-checking.
func combineSamples(samples []clockSample, target int) (time.Duration, string, error) {
	slices.SortStableFunc(samples, func(a, b clockSample) int { return cmp.Compare(a.distance, b.distance) })
	var point time.Duration
	bestCount := 0
	for _, candidate := range samples {
		lower := candidate.offset - candidate.distance
		count := 0
		for _, sample := range samples {
			if lower >= sample.offset-sample.distance && lower <= sample.offset+sample.distance {
				count++
			}
		}
		if count > bestCount {
			bestCount, point = count, lower
		}
	}
	if bestCount <= len(samples)/2 {
		return 0, "", errors.New("NTP servers disagree on clock offset")
	}

	var selected []clockSample
	for _, sample := range samples {
		if point >= sample.offset-sample.distance && point <= sample.offset+sample.distance {
			selected = append(selected, sample)
		}
	}
	best := selected[0]
	lower, upper := best.offset-best.distance, best.offset+best.distance
	for _, sample := range selected {
		lower = max(lower, sample.offset-sample.distance)
		upper = min(upper, sample.offset+sample.distance)
	}
	var weightedDelta, weightSum float64
	for _, sample := range selected[:min(target, len(selected))] {
		ratio := float64(best.distance) / float64(sample.distance)
		weight := ratio * ratio
		// Subtract an anchor before conversion: monotonic-clock offsets can
		// span decades, which would lose sub-microsecond precision as float64.
		weightedDelta += float64(sample.offset-best.offset) * weight
		weightSum += weight
	}
	offset := best.offset + time.Duration(math.Round(weightedDelta/weightSum))
	return min(max(offset, lower), upper), best.server, nil
}
