package cache

import "time"

func (c *AnyCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache = nil
	c.createdAt = time.Time{}
}

func (c *KvCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache = map[any]any{}
	c.createdAt = time.Time{}
}
