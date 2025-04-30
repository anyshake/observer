package explorer

import "time"

func (s *DeviceStatus) SetStartedAt(startedAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.startedAt = startedAt
}

func (s *DeviceStatus) GetStartedAt() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.startedAt
}

func (s *DeviceStatus) SetUpdatedAt(updatedAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.updatedAt = updatedAt
}

func (s *DeviceStatus) GetUpdatedAt() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.updatedAt
}

func (s *DeviceStatus) IncrementFrames() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.frames++
}

func (s *DeviceStatus) GetFrames() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.frames
}

func (s *DeviceStatus) IncrementErrors() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.errors++
}

func (s *DeviceStatus) GetErrors() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.errors
}

func (s *DeviceStatus) IncrementMessages() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messages++
}

func (s *DeviceStatus) GetMessages() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.messages
}
