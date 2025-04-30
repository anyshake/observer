package message

import (
	"fmt"
)

func (s *Bus[T]) Subscribe(clientId string, handler T) error {
	if _, ok := s.subscribers.Get(clientId); ok {
		return fmt.Errorf("client %s is already subscribed", clientId)
	}
	s.subscribers.Set(clientId, handler)

	return s.messageBus.Subscribe(s.topicName, handler)
}
