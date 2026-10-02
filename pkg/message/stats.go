package message

func (s *Bus[T]) GetStats() Stats {
	subscriberCount := 0
	if snapshot := s.snapshot.Load(); snapshot != nil {
		subscriberCount = len(snapshot.subscribers)
	}

	return Stats{
		Published:    s.published.Load(),
		Delivered:    s.delivered.Load(),
		Dropped:      s.dropped.Load(),
		Disconnected: s.disconnected.Load(),
		Subscribers:  subscriberCount,
	}
}
