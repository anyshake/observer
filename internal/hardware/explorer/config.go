package explorer

import "time"

// Snapshot returns an independent copy of the configuration with a fresh mutex.
func (c *DeviceConfig) Snapshot() *DeviceConfig {
	c.mu.Lock()
	defer c.mu.Unlock()

	return &DeviceConfig{
		packetInterval: c.packetInterval,
		channelCodes:   append([]string(nil), c.channelCodes...),
		sampleRate:     c.sampleRate,
		gnssEnabled:    c.gnssEnabled,
		model:          c.model,
		protocol:       c.protocol,
	}
}

func (c *DeviceConfig) SetPacketInterval(packetInterval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.packetInterval = packetInterval
}

func (c *DeviceConfig) GetPacketInterval() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.packetInterval
}

func (c *DeviceConfig) SetSampleRate(sampleRate int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.sampleRate = sampleRate
}

func (c *DeviceConfig) GetSampleRate() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.sampleRate
}

func (s *DeviceConfig) SetChannelCodes(channels []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.channelCodes = channels
}

func (s *DeviceConfig) GetChannelCodes() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.channelCodes
}

func (c *DeviceConfig) SetGnssAvailability(gnssEnabled bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.gnssEnabled = gnssEnabled
}

func (c *DeviceConfig) GetGnssAvailability() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.gnssEnabled
}

func (s *DeviceConfig) SetModel(model string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.model = model
}

func (s *DeviceConfig) GetModel() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.model
}

func (s *DeviceConfig) SetProtocol(protocol string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.protocol = protocol
}

func (s *DeviceConfig) GetProtocol() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.protocol
}
