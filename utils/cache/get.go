package cache

func (a *AnyCache) Get() any {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.cache
}

func (c *KvCache) Get(key any) (any, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	data, ok := c.cache[key]

	return data, ok
}
