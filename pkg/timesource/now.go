package timesource

import (
	"time"
)

func (g *Source) Now() time.Time {
	g.mu.RLock()
	tm := g.timeFunc()
	local := g.localTime
	ref := g.refTime
	ppm := g.driftPPM
	g.mu.RUnlock()

	elapsed := tm.Sub(local)
	drift := time.Duration(elapsed.Seconds() * ppm * 1e-6 * float64(time.Second))

	return ref.Add(elapsed + drift).UTC()
}
