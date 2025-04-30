package cache

import "time"

func (c *AnyCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache = nil
	c.createdAt = time.Unix(0, 0)
}

func (c *KvCache[T]) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache = map[any]T{}
	c.createdAt = time.Unix(0, 0)
}
