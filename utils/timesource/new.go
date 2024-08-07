package timesource

import "time"

func New(baseTime, refTime time.Time) Source {
	return Source{
		LocalBaseTime: baseTime,
		ReferenceTime: refTime,
	}
}
