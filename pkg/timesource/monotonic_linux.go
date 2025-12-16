package timesource

import (
	"time"

	"golang.org/x/sys/unix"
)

func Monotonic() time.Duration {
	var ts unix.Timespec
	unix.ClockGettime(unix.CLOCK_MONOTONIC_RAW, &ts)
	return time.Duration(ts.Nano())
}

func MonotonicNow() time.Time {
	y2k := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	return y2k.Add(Monotonic()) // starting from Jan 1, 2000
}
