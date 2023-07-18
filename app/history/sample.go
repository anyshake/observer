package history

import (
	"fmt"
	"time"

	"com.geophone.observer/handler"
	t "com.geophone.observer/utils/time"
)

func getSampleRate(data []handler.Geophone, channel string) (int, error) {
	var (
		length   = 0
		lastTime = time.UnixMilli(data[0].TS)
	)

	for _, v := range data {
		currentTime := time.UnixMilli(v.TS)
		if t.Diff(currentTime, lastTime) > THRESHOLD {
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
