package winston

import (
	"context"
	"fmt"
)

func (s *WinstonServiceImpl) Restart(ctx context.Context) error {
	defer s.status.IncrementRestarts()

	if err := s.Stop(ctx); err != nil {
		return fmt.Errorf("restart failed when stopping service: %w", err)
	}

	if err := s.Init(ctx); err != nil {
		return fmt.Errorf("restart failed when initializing service: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	if err := s.Start(ctx); err != nil {
		return fmt.Errorf("restart failed when starting service: %w", err)
	}

	return nil
}
