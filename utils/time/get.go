package time

import "time"

func Get(offset float64) int64 {
	current := time.Now().UTC().UnixMilli()

	return int64(offset*float64(time.Second.Milliseconds())) + current
}
