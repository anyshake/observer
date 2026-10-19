package archiver

import "context"

func (s *ArchiverServiceImpl) IsEnabled(ctx context.Context) bool {
	enable, err := (&archiverConfigEnabledImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return false
	}
	if en, ok := enable.(bool); !ok || !en {
		return false
	}
	return true
}
