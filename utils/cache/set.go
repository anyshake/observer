package cache

import "time"

func (c *BytesCache) Set(data []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache = data
	c.createdAt = time.Now()
}
