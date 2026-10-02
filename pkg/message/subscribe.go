package message

import (
	"fmt"
	"runtime/debug"
	"time"
)

func (s *Bus[T]) Subscribe(clientID string, options SubscriptionOptions, handler func(T)) error {
	if handler == nil {
		return ErrInvalidHandler
	}
	options = normalizeSubscriptionOptions(options)
	var errorQueue chan error
	if options.OnError != nil {
		errorQueue = make(chan error, 1)
	}
	subscriber := &subscriber[T]{
		handler:    handler,
		options:    options,
		queue:      make(chan T, options.BufferSize),
		errorQueue: errorQueue,
		done:       make(chan struct{}),
	}

	s.mu.Lock()
	if s.closed.Load() {
		s.mu.Unlock()
		return ErrBusClosed
	}
	if _, ok := s.subscribers[clientID]; ok {
		s.mu.Unlock()
		return fmt.Errorf("client %s is already subscribed", clientID)
	}
	s.subscribers[clientID] = subscriber
	s.refreshSnapshotLocked()
	s.mu.Unlock()

	go s.runSubscriber(clientID, subscriber)
	if options.OnError != nil {
		go subscriber.runErrorHandler()
	}
	return nil
}

func (s *Bus[T]) refreshSnapshotLocked() {
	subscribers := make([]subscriberRef[T], 0, len(s.subscribers))
	for clientID, subscriber := range s.subscribers {
		subscribers = append(subscribers, subscriberRef[T]{
			clientID:   clientID,
			subscriber: subscriber,
		})
	}
	s.snapshot.Store(&subscriberSnapshot[T]{subscribers: subscribers})
}

func normalizeSubscriptionOptions(options SubscriptionOptions) SubscriptionOptions {
	if options.BufferSize < 1 {
		options.BufferSize = 1
	}
	if options.Overflow == OverflowBlock && options.EnqueueTimeout <= 0 {
		options.EnqueueTimeout = 250 * time.Millisecond
	}
	return options
}

func (s *Bus[T]) runSubscriber(clientID string, subscriber *subscriber[T]) {
	for {
		select {
		case <-subscriber.done:
			return
		default:
		}

		select {
		case <-subscriber.done:
			return
		case event := <-subscriber.queue:
			if panicValue := callHandler(subscriber.handler, event); panicValue != nil {
				err := fmt.Errorf("%w for client %s: %v\n%s", ErrSubscriberPanicked, clientID, panicValue, debug.Stack())
				if s.disconnect(clientID, subscriber, err) {
					s.disconnected.Add(1)
				}
				return
			}
		}
	}
}

func callHandler[T any](handler func(T), event T) (panicValue any) {
	defer func() {
		panicValue = recover()
	}()
	handler(event)
	return nil
}

func (s *subscriber[T]) stop(err error) {
	s.stopOnce.Do(func() {
		if err != nil {
			s.replacePendingError(err)
		}
		close(s.done)
	})
}

func (s *subscriber[T]) replacePendingError(err error) {
	if s.options.OnError == nil {
		return
	}
	select {
	case <-s.errorQueue:
	default:
	}
	s.notifyError(err)
}

func (s *subscriber[T]) notifyError(err error) {
	if s.options.OnError == nil {
		return
	}
	select {
	case s.errorQueue <- err:
	default:
	}
}

func (s *subscriber[T]) runErrorHandler() {
	for {
		select {
		case err := <-s.errorQueue:
			s.callErrorHandler(err)
		case <-s.done:
			select {
			case err := <-s.errorQueue:
				s.callErrorHandler(err)
			default:
			}
			return
		}
	}
}

func (s *subscriber[T]) callErrorHandler(err error) {
	defer func() {
		_ = recover()
	}()
	s.options.OnError(err)
}
