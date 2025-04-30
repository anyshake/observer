package quakesense

import "fmt"

func (s *QuakeSenseServiceImpl) Restart() error {
	defer s.status.IncrementRestarts()

	if err := s.Stop(); err != nil {
		return fmt.Errorf("restart failed when stopping service: %w", err)
	}

	if err := s.Init(); err != nil {
		return fmt.Errorf("restart failed when initializing service: %w", err)
	}

	if err := s.Start(); err != nil {
		return fmt.Errorf("restart failed when starting service: %w", err)
	}

	return nil
}
