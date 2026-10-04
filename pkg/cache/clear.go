package cache

import "time"

func (c *GenericCache[T]) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache = nil
	c.createdAt = time.Time{}
}

func (c *KvCache[T]) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache = map[any]T{}
	c.createdAt = time.Time{}
}
