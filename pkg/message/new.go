package message

func NewBus[T any](topicName string) *Bus[T] {
	bus := &Bus[T]{
		topicName:   topicName,
		subscribers: make(map[string]*subscriber[T]),
	}
	bus.snapshot.Store(&subscriberSnapshot[T]{})
	return bus
}
