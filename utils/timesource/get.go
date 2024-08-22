package timesource

import (
	"errors"
	"time"
)

func (g *Source) Get() (time.Time, error) {
	g.rwMutex.RLock()
	defer g.rwMutex.RUnlock()

	if g.LocalBaseTime.IsZero() || g.ReferenceTime.IsZero() {
		return time.Now().UTC(), errors.New("empty BaseTime or RefTime is not allowed")
	}

	elapsed := time.Since(g.LocalBaseTime.UTC())
	return g.ReferenceTime.Add(elapsed).UTC(), nil
}
