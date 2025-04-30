package seedlink

func (s *SeedLinkServiceImpl) GetDescription() string {
	return "A simple built-in SeedLink server that accepts third-party clients (e.g. Swarm) to connect to this station. Note that this implementation is not a standard SeedLink protocol and may have compatibility issues with some clients."
}
