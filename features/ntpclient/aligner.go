package ntpclient

import "time"

func AlignTime(offset float64) int64 {
	current := time.Now().UnixMilli()

	return int64(offset*float64(time.Second.Milliseconds())) + current
}
