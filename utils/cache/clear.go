package cache

import "time"

func (c *BytesCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache = nil
	c.createdAt = time.Time{}
}
