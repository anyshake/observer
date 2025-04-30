package message

func (s *Bus[T]) GetTopicName() string {
	return s.topicName
}
