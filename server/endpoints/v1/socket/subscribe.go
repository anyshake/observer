package socket

import (
	"errors"

	"github.com/anyshake/observer/drivers/explorer"
)

func (s *Socket) subscribe(clientId string, handler explorer.ExplorerEventHandler) error {
	if _, ok := s.subscribers.Get(clientId); ok {
		return errors.New("this client has already subscribed")
	}
	s.subscribers.Set(clientId, handler)
	s.messageBus.Subscribe(s.GetApiName(), handler)
	return nil
}
