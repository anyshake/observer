package explorer

import "math"

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

func (c *ExplorerConfig) SetDeviceInfo(deviceInfo uint32) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.deviceInfo = deviceInfo
}

func (c *ExplorerConfig) GetDeviceInfo() (deviceInfo, deviceId uint32) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Leagcy mode is always 0xFFFFFFFF
	if c.deviceInfo == math.MaxUint32 {
		return c.deviceInfo, c.deviceInfo
	}

	// 0x7FFFFFFF represents a device with no assigned serial number
	deviceId = c.deviceInfo & 0x7FFFFFFF
	if deviceId == 0x7FFFFFFF {
		deviceId = math.MaxUint32
	}

	return c.deviceInfo, deviceId
}

func (c *ExplorerConfig) SetLatitude(latitude float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.deviceInfo&0x80000000 == 0 {
		c.fallbackLatitude = latitude
	} else {
		c.gnssLatitude = latitude
	}
}

func (c *ExplorerConfig) GetLatitude() float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.deviceInfo&0x80000000 == 0 {
		return c.fallbackLatitude
	}

	return c.gnssLatitude
}

func (c *ExplorerConfig) SetLongitude(longitude float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.deviceInfo&0x80000000 == 0 {
		c.fallbackLongitude = longitude
	} else {
		c.gnssLongitude = longitude
	}
}

func (c *ExplorerConfig) GetLongitude() float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.deviceInfo&0x80000000 == 0 {
		return c.fallbackLongitude
	}

	return c.gnssLongitude
}

func (c *ExplorerConfig) SetElevation(elevation float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.deviceInfo&0x80000000 == 0 {
		c.fallbackElevation = elevation
	} else {
		c.gnssElevation = elevation
	}
}

func (c *ExplorerConfig) GetElevation() float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.deviceInfo&0x80000000 == 0 {
		return c.fallbackElevation
	}

	return c.gnssElevation
}
