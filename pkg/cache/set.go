package cache

import "time"

func (a *AnyCache) Set(data any) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.cache = data
	a.createdAt = time.Now()
}

func (c *KvCache[T]) Set(key any, data T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache[key] = data
	c.createdAt = time.Now()
}
