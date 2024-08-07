package timesource

import "time"

type Source struct {
	LocalBaseTime time.Time
	ReferenceTime time.Time
}
