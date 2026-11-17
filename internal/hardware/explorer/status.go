package explorer

import "time"

// Snapshot returns an independent copy of the status with a fresh mutex.
func (s *DeviceStatus) Snapshot() *DeviceStatus {
	s.mu.Lock()
	defer s.mu.Unlock()

	return &DeviceStatus{
		startedAt:  s.startedAt,
		updatedAt:  s.updatedAt,
		frames:     s.frames,
		errors:     s.errors,
		messages:   s.messages,
		gnssStatus: s.gnssStatus.clone(),
	}
}

func (s *GnssStatus) clone() *GnssStatus {
	if s == nil {
		return nil
	}
	snapshot := *s
	if s.HDOP != nil {
		hdop := *s.HDOP
		snapshot.HDOP = &hdop
	}
	if s.PPSAgeSeconds != nil {
		age := *s.PPSAgeSeconds
		snapshot.PPSAgeSeconds = &age
	}
	return &snapshot
}

func (s *DeviceStatus) SetGnssStatus(status *GnssStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.gnssStatus = status.clone()
}

// GetGnssStatus returns nil until firmware advertising variable 6 reports it.
func (s *DeviceStatus) GetGnssStatus() *GnssStatus {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.gnssStatus.clone()
}

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
