package message

func (s *Bus[T]) Close() {
	s.mu.Lock()
	if s.closed.Load() {
		s.mu.Unlock()
		return
	}
	s.closed.Store(true)
	subscribers := make([]*subscriber[T], 0, len(s.subscribers))
	for clientID, subscriber := range s.subscribers {
		subscribers = append(subscribers, subscriber)
		delete(s.subscribers, clientID)
	}
	s.refreshSnapshotLocked()
	s.mu.Unlock()

	for _, subscriber := range subscribers {
		subscriber.stop(ErrBusClosed)
	}
}
