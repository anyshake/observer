package cache

import "time"

func (c *AnyCache) Valid() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return time.Since(c.createdAt) < c.ttl
}
