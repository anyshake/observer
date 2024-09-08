package forwarder

import (
	"errors"

	"github.com/anyshake/observer/drivers/explorer"
)

func (s *ForwarderService) subscribe(clientId string, handler explorer.ExplorerEventHandler) error {
	if _, ok := s.subscribers.Get(clientId); ok {
		return errors.New("this client has already subscribed")
	}
	s.subscribers.Set(clientId, handler)
	s.messageBus.Subscribe(s.GetServiceName(), handler)
	return nil
}
