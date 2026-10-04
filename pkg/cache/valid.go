package cache

import "time"

func (c *GenericCache[T]) Valid() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return !c.createdAt.IsZero() && time.Since(c.createdAt) < c.ttl
}

func (c *KvCache[T]) Valid() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return !c.createdAt.IsZero() && time.Since(c.createdAt) < c.ttl
}
