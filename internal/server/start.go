package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (s *httpServer) Start() error {
	if s.debug {
		go func() {
			for {
				warnText := "PLEASE DISABLE DEBUG MODE IN PRODUCTION ENVIRONMENT TO AVOID SECURITY RISKS"
				border := strings.Repeat("=", len(warnText)+6)
				s.log.Warnf("\n%s\n|| %s ||\n%s", border, warnText, border)
				time.Sleep(time.Minute)
			}
		}()
	}

	s.log.Infof("starting http server at %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	return nil
}
