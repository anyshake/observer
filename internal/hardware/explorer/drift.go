package explorer

import (
	"math"
	"time"

	"github.com/anyshake/observer/pkg/ringbuf"
)

type clockDrift struct {
	offset     time.Duration
	measuredAt time.Time
}

func getLongTermClockDriftPPM(buf *ringbuf.Buffer[clockDrift], ppmWindow time.Duration) float64 {
	vals := buf.Values()
	if len(vals) <= 1 {
		return 0
	}

	now := vals[len(vals)-1].measuredAt
	cutoff := now.Add(-ppmWindow)

	start := 0
	for start < len(vals) && vals[start].measuredAt.Before(cutoff) {
		start++
	}

	if len(vals)-start <= 1 {
		return 0
	}

	if start > 0 {
		buf.Reset()
		buf.Push(vals[start:]...)
		vals = buf.Values()
	}

	oldest := vals[0]
	newest := vals[len(vals)-1]

	dt := newest.measuredAt.Sub(oldest.measuredAt).Seconds()
	if dt <= 0 {
		return 0
	}

	dOffset := (newest.offset - oldest.offset).Seconds()
	ppm := dOffset / dt * 1e6

	if math.Abs(ppm) >= 50 {
		buf.Reset()
		ppm = 0
	}

	return ppm
}
