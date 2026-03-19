package mdns_discovery

func (s *DiscoveryServiceImpl) IsEnabled() bool {
	enable, err := (&discoveryConfigEnabledImpl{}).Get(s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
