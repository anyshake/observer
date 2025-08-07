package timesource

import (
	"time"
)

func (g *Source) Now() time.Time {
	g.mu.Lock()
	defer g.mu.Unlock()

	elapsed := time.Since(g.localTime)
	return g.refTime.Add(elapsed).UTC()
}
