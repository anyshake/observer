package server

import (
	"context"
	"fmt"
	"time"
)

func (s *httpServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}

	return nil
}
