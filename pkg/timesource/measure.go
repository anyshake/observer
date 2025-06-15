package timesource

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/beevik/ntp"
)

func measureStableOffset(ntpAddress string, retries int, timeout time.Duration, samples int) (time.Duration, error) {
	var rawOffsets []time.Duration

	for i := 0; i < samples; i++ {
		var (
			res *ntp.Response
			err error
		)

		for attempt := 0; attempt <= retries; attempt++ {
			res, err = ntp.QueryWithOptions(ntpAddress, ntp.QueryOptions{Timeout: timeout})
			if err == nil {
				rawOffsets = append(rawOffsets, res.ClockOffset)
				break
			}
		}

		if err != nil {
			return 0, fmt.Errorf("failed to query NTP server after %d retries (iteration %d): %w", retries, i+1, err)
		}

		time.Sleep(100 * time.Millisecond)
	}

	if len(rawOffsets) == 0 {
		return 0, errors.New("no valid NTP responses")
	}

	var sum time.Duration
	for _, offset := range rawOffsets {
		sum += offset
	}
	mean := sum / time.Duration(len(rawOffsets))

	var variance float64
	for _, offset := range rawOffsets {
		diff := float64(offset - mean)
		variance += diff * diff
	}
	stddev := time.Duration(math.Sqrt(variance / float64(len(rawOffsets))))

	var filtered []time.Duration
	for _, offset := range rawOffsets {
		if absDuration(offset-mean) <= 2*stddev {
			filtered = append(filtered, offset)
		}
	}

	if len(filtered) == 0 {
		return 0, errors.New("all NTP offsets were filtered as outliers")
	}

	sum = 0
	for _, offset := range filtered {
		sum += offset
	}
	return sum / time.Duration(len(filtered)), nil
}

func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
