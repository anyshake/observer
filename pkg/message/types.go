package message

import (
	"github.com/alphadose/haxmap"
	messagebus "github.com/vardius/message-bus"
)

type Bus[T any] struct {
	topicName   string
	messageBus  messagebus.MessageBus
	subscribers *haxmap.Map[string, T]
}
