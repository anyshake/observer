package seedlink

import "context"

func (s *SeedLinkServiceImpl) IsEnabled(ctx context.Context) bool {
	enable, err := (&seedlinkConfigEnabledImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
