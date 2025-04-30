package message

import (
	"github.com/alphadose/haxmap"
	messagebus "github.com/vardius/message-bus"
)

func NewBus[T any](topicName string, size int) Bus[T] {
	return Bus[T]{
		messageBus:  messagebus.New(size),
		subscribers: haxmap.New[string, T](),
	}
}
