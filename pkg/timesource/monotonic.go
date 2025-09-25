package timesource

import (
	"time"
	_ "unsafe"
)

//go:noescape
//go:linkname nanotime runtime.nanotime
func nanotime() int64

func Monotonic() time.Duration {
	return time.Duration(nanotime())
}

func MonotonicNow() time.Time {
	y2k := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	return y2k.Add(Monotonic()) // starting from Jan 1, 2000
}
