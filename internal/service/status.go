package service

import "time"

func (s *Status) SetIsRunning(isRunning bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.isRunning = isRunning
}

func (s *Status) GetIsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.isRunning
}

func (s *Status) IncrementRestarts() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.restarts++
}

func (s *Status) SetRestarts(restarts int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.restarts = restarts
}

func (s *Status) GetRestarts() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.restarts
}

func (s *Status) SetStartedAt(startedAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.startedAt = startedAt
}

func (s *Status) GetStartedAt() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.startedAt
}

func (s *Status) SetStoppedAt(stoppedAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stoppedAt = stoppedAt
}

func (s *Status) GetStoppedAt() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.stoppedAt
}

func (s *Status) SetUpdatedAt(updatedAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.updatedAt = updatedAt
}

func (s *Status) GetUpdatedAt() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.updatedAt
}
