package mdns_discovery

import "context"

func (s *DiscoveryServiceImpl) IsEnabled(ctx context.Context) bool {
	enable, err := (&discoveryConfigEnabledImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
