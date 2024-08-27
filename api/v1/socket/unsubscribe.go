package socket

import (
	"errors"
)

func (s *Socket) unsubscribe(clientId string) error {
	fn, ok := s.subscribers.Get(clientId)
	if !ok {
		return errors.New("this client has not subscribed")
	}

	s.subscribers.Del(clientId)
	return s.messageBus.Unsubscribe(s.GetApiName(), fn)
}
