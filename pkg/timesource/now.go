package timesource

import (
	"time"
)

func (g *Source) Now() time.Time {
	g.mu.RLock()
	defer g.mu.RUnlock()

	elapsed := time.Since(g.localTime)
	return g.refTime.Add(elapsed).UTC()
}
