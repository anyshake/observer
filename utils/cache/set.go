package cache

import "time"

func (a *AnyCache) Set(data any) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.cache = data
	a.createdAt = time.Now()
}
