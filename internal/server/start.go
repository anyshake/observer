package server

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

func (s *HttpServer) Listen() error {
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen http server at %s: %w", s.server.Addr, err)
	}
	s.listener = listener
	return nil
}

func (s *HttpServer) Start() error {
	if s.listener == nil {
		return errors.New("http server is not listening")
	}

	if s.debug {
		go s.warnDebugMode()
	}

	s.log.Infof("starting http server at %s", s.server.Addr)
	if err := s.server.Serve(s.listener); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	return nil
}

func (s *HttpServer) warnDebugMode() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	warnText := "PLEASE DISABLE DEBUG MODE IN PRODUCTION ENVIRONMENT TO AVOID SECURITY RISKS"
	border := strings.Repeat("=", len(warnText)+6)
	for {
		s.log.Warnf("\n%s\n|| %s ||\n%s", border, warnText, border)
		select {
		case <-s.shutdownCtx.Done():
			return
		case <-ticker.C:
		}
	}
}
