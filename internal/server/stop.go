package server

import (
	"context"
	"errors"
	"fmt"
	"net"
)

func (s *HttpServer) Stop(ctx context.Context) error {
	s.shutdownCancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}
	if s.listener != nil {
		if err := s.listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			return fmt.Errorf("failed to close http listener: %w", err)
		}
	}
	if s.closeSocket != nil {
		if err := s.closeSocket(ctx); err != nil {
			return fmt.Errorf("failed to close websocket stream: %w", err)
		}
	}

	return nil
}
