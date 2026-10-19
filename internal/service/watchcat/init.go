package watchcat

import "context"

func (s *WatchCatServiceImpl) Init(ctx context.Context) error {
	return ctx.Err()
}
