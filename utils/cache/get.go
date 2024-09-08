package cache

func (c *BytesCache) Get() []byte {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.cache
}
