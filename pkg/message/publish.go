package message

func (s *Bus[T]) Publish(message ...any) {
	s.messageBus.Publish(s.topicName, message...)
}
