package forwarder

import (
	"errors"
)

func (s *ForwarderService) unsubscribe(clientId string) error {
	fn, ok := s.subscribers.Get(clientId)
	if !ok {
		return errors.New("this client has not subscribed")
	}
	s.messageBus.Unsubscribe(s.GetServiceName(), fn)
	s.subscribers.Remove(clientId)
	return nil
}
