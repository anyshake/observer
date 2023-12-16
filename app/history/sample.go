package history

import (
	"fmt"
	"time"

	"github.com/anyshake/observer/publisher"
	"github.com/anyshake/observer/utils/duration"
)

func getSampleRate(data []publisher.Geophone, channel string) (int, error) {
	var (
		sampleRateSum float64
		lastTime      = time.UnixMilli(data[0].TS)
	)

	for i := 1; i < len(data); i++ {
		var (
			currentTime = time.UnixMilli(data[i].TS)
			timeDiff    = duration.Difference(currentTime, lastTime)
		)
		if timeDiff > THRESHOLD {
			err := fmt.Errorf("uneven gaps between the data")
			return 0, err
		}

		switch channel {
		case "EHZ":
			sampleRateSum += float64(len(data[i].EHZ)) / timeDiff.Seconds()
		case "EHE":
			sampleRateSum += float64(len(data[i].EHE)) / timeDiff.Seconds()
		case "EHN":
			sampleRateSum += float64(len(data[i].EHN)) / timeDiff.Seconds()
		}

		lastTime = currentTime
	}

	return int(sampleRateSum) / len(data), nil
}
