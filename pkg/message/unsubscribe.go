package message

import (
	"fmt"
)

func (s *Bus[T]) Unsubscribe(clientId string) error {
	s.mu.Lock()
	subscriber, ok := s.subscribers[clientId]
	if !ok {
		s.mu.Unlock()
		return fmt.Errorf("client %s is not subscribed", clientId)
	}
	delete(s.subscribers, clientId)
	s.refreshSnapshotLocked()
	s.mu.Unlock()

	subscriber.stop(nil)
	return nil
}

func (s *Bus[T]) disconnect(clientID string, expected *subscriber[T], err error) bool {
	s.mu.Lock()
	subscriber, ok := s.subscribers[clientID]
	if !ok || subscriber != expected {
		s.mu.Unlock()
		return false
	}
	delete(s.subscribers, clientID)
	s.refreshSnapshotLocked()
	s.mu.Unlock()

	subscriber.stop(err)
	return true
}
