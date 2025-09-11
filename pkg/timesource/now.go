package timesource

import (
	"time"
)

func (g *Source) Now() time.Time {
	elapsed := time.Since(g.localTime)

	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.refTime.Add(elapsed).UTC()
}
