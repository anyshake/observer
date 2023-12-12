package history

import (
	"fmt"
	"time"

	"github.com/anyshake/observer/publisher"
	"github.com/anyshake/observer/utils/duration"
)

func getSampleRate(data []publisher.Geophone, channel string) (int, error) {
	var (
		length   = 0
		lastTime = time.UnixMilli(data[0].TS)
	)

	for _, v := range data {
		currentTime := time.UnixMilli(v.TS)
		if duration.Difference(currentTime, lastTime) > THRESHOLD {
			err := fmt.Errorf("uneven gaps between the data")
			return 0, err
		}

		switch channel {
		case "EHZ":
			length += len(v.EHZ)
		case "EHE":
			length += len(v.EHE)
		case "EHN":
			length += len(v.EHN)
		}

		lastTime = time.UnixMilli(v.TS)
	}

	return length / len(data), nil
}
