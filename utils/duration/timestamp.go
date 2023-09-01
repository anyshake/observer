package duration

import "time"

func Timestamp(offset float64) (time.Time, int64) {
	offsetDuration := time.Duration(offset * float64(time.Second))
	currentTime := time.Now().UTC()

	result := currentTime.Add(offsetDuration)
	return result, result.UnixMilli()
}
