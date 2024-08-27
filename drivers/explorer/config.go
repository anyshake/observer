package explorer

func (c *ExplorerConfig) SetLegacyMode(legacyMode bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.legacyMode = legacyMode
}

func (c *ExplorerConfig) GetLegacyMode() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.legacyMode
}

func (c *ExplorerConfig) SetDeviceId(deviceId uint32) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.deviceId = deviceId
}

func (c *ExplorerConfig) GetDeviceId() uint32 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.deviceId
}

func (c *ExplorerConfig) SetLatitude(latitude float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.latitude = latitude
}

func (c *ExplorerConfig) GetLatitude() float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.latitude
}

func (c *ExplorerConfig) SetLongitude(longitude float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.longitude = longitude
}

func (c *ExplorerConfig) GetLongitude() float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.longitude
}

func (c *ExplorerConfig) SetElevation(elevation float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.elevation = elevation
}

func (c *ExplorerConfig) GetElevation() float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.elevation
}
