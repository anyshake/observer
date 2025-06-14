package timesource

import (
	"time"
)

func (g *Source) Get() time.Time {
	g.rwMutex.RLock()
	defer g.rwMutex.RUnlock()

	elapsed := time.Since(g.LocalBaseTime)
	return g.ReferenceTime.Add(elapsed).UTC()
}
