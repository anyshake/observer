package cache

func (a *AnyCache) Get() any {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.cache
}
