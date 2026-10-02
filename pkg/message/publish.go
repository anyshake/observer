package message

import "time"

type enqueueResult uint8

const (
	enqueueDelivered enqueueResult = iota
	enqueueDropped
	enqueueReplaced
	enqueueDisconnected
)

func (s *Bus[T]) Publish(event T) PublishResult {
	s.published.Add(1)

	if s.closed.Load() {
		return PublishResult{}
	}
	snapshot := s.snapshot.Load()
	if snapshot == nil {
		return PublishResult{}
	}

	result := PublishResult{}
	for _, ref := range snapshot.subscribers {
		switch ref.subscriber.enqueue(event) {
		case enqueueDelivered:
			result.Delivered++
			s.delivered.Add(1)
		case enqueueDropped:
			result.Dropped++
			s.dropped.Add(1)
		case enqueueReplaced:
			result.Delivered++
			result.Dropped++
			s.delivered.Add(1)
			s.dropped.Add(1)
		case enqueueDisconnected:
			if s.disconnect(ref.clientID, ref.subscriber, ErrSubscriberTooSlow) {
				result.Disconnected++
				s.disconnected.Add(1)
			}
		}
	}

	return result
}

func (s *subscriber[T]) enqueue(event T) enqueueResult {
	select {
	case <-s.done:
		return enqueueDropped
	default:
	}

	switch s.options.Overflow {
	case OverflowBlock:
		select {
		case s.queue <- event:
			return enqueueDelivered
		default:
		}

		timer := time.NewTimer(s.options.EnqueueTimeout)
		defer timer.Stop()
		select {
		case s.queue <- event:
			return enqueueDelivered
		case <-s.done:
			return enqueueDropped
		case <-timer.C:
			s.notifyError(ErrPublishTimeout)
			return enqueueDropped
		}
	case OverflowDropOldest:
		select {
		case s.queue <- event:
			return enqueueDelivered
		default:
		}

		select {
		case <-s.queue:
		default:
		}
		select {
		case s.queue <- event:
			return enqueueReplaced
		case <-s.done:
			return enqueueDropped
		default:
			return enqueueDropped
		}
	case OverflowDisconnect:
		select {
		case s.queue <- event:
			return enqueueDelivered
		case <-s.done:
			return enqueueDropped
		default:
			return enqueueDisconnected
		}
	default:
		select {
		case s.queue <- event:
			return enqueueDelivered
		case <-s.done:
			return enqueueDropped
		default:
			return enqueueDropped
		}
	}
}
