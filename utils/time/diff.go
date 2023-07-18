package time

import "time"

func Diff(a, b time.Time) time.Duration {
	if a.After(b) {
		return a.Sub(b)
	}

	return b.Sub(a)
}
