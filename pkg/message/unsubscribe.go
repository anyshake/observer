package message

import (
	"fmt"
)

func (s *Bus[T]) Unsubscribe(clientId string) error {
	fn, ok := s.subscribers.Get(clientId)
	if !ok {
		return fmt.Errorf("client %s is not subscribed", clientId)
	}
	s.subscribers.Del(clientId)

	return s.messageBus.Unsubscribe(s.topicName, fn)
}
