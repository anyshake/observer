package forwarder

import (
	"errors"
)

func (s *ForwarderService) unsubscribe(clientId string) error {
	fn, ok := s.subscribers.Get(clientId)
	if !ok {
		return errors.New("this client has not subscribed")
	}
	s.subscribers.Del(clientId)
	return s.messageBus.Unsubscribe(s.GetServiceName(), fn)
}
