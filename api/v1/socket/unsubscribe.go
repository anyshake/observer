package socket

import (
	"errors"
)

func (s *Socket) Unsubscribe(clientId string) error {
	fn, ok := s.subscribers.Get(clientId)
	if !ok {
		return errors.New("this client has not subscribed")
	}
	s.messageBus.Unsubscribe(s.GetApiName(), fn)
	s.subscribers.Remove(clientId)
	return nil
}
