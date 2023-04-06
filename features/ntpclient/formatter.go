package ntpclient

import "time"

func FormatTime(timestamp int64, format string) string {
	return time.Unix(0, timestamp*int64(time.Millisecond)).Format(format)
}
