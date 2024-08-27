package seedlink

import (
	"errors"
	"fmt"

	"github.com/alphadose/haxmap"
	"github.com/anyshake/observer/drivers/explorer"
	"github.com/bclswl0827/slgo/handlers"
	messagebus "github.com/vardius/message-bus"
)

type consumer struct {
	channelPrefix string
	serviceName   string
	messageBus    messagebus.MessageBus // An independent message bus for the socket module
	subscribers   *haxmap.Map[string, func(data *explorer.ExplorerData)]
}

func (c *consumer) Subscribe(clientId string, channels []string, eventHandler func(handlers.SeedLinkDataPacket)) error {
	if _, ok := c.subscribers.Get(clientId); ok {
		return errors.New("this client has already subscribed")
	}
	handler := func(data *explorer.ExplorerData) {
		for _, channel := range channels {
			switch channel {
			case fmt.Sprintf("%sZ", c.channelPrefix):
				eventHandler(handlers.SeedLinkDataPacket{
					Timestamp:  data.Timestamp,
					SampleRate: data.SampleRate,
					Channel:    channel,
					DataArr:    data.Z_Axis,
				})
			case fmt.Sprintf("%sE", c.channelPrefix):
				eventHandler(handlers.SeedLinkDataPacket{
					Timestamp:  data.Timestamp,
					SampleRate: data.SampleRate,
					Channel:    channel,
					DataArr:    data.E_Axis,
				})
			case fmt.Sprintf("%sN", c.channelPrefix):
				eventHandler(handlers.SeedLinkDataPacket{
					Timestamp:  data.Timestamp,
					SampleRate: data.SampleRate,
					Channel:    channel,
					DataArr:    data.N_Axis,
				})
			}
		}
	}
	c.subscribers.Set(clientId, handler)
	c.messageBus.Subscribe(c.serviceName, handler)
	return nil
}

func (c *consumer) Unsubscribe(clientId string) error {
	fn, ok := c.subscribers.Get(clientId)
	if !ok {
		return errors.New("this client has not subscribed")
	}

	c.subscribers.Del(clientId)
	return c.messageBus.Unsubscribe(c.serviceName, fn)
}
