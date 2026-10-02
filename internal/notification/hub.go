package notification

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Hub struct {
	mu                   sync.RWMutex
	subscriberBufferSize int
	subscribers          map[chan Event]struct{}
}

func NewHub(subscriberBufferSize int) *Hub {
	if subscriberBufferSize < 1 {
		subscriberBufferSize = 1
	}

	return &Hub{
		subscriberBufferSize: subscriberBufferSize,
		subscribers:          make(map[chan Event]struct{}),
	}
}

func (h *Hub) Publish(event Event) {
	if event.ID == "" {
		event.ID = uuid.NewString()
	}
	if event.OccurredAt == 0 {
		event.OccurredAt = time.Now().UnixMilli()
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for subscriber := range h.subscribers {
		select {
		case subscriber <- event:
		default:
		}
	}
}

func (h *Hub) Subscribe() (<-chan Event, func()) {
	subscriber := make(chan Event, h.subscriberBufferSize)
	h.mu.Lock()
	h.subscribers[subscriber] = struct{}{}
	h.mu.Unlock()

	var unsubscribeOnce sync.Once
	unsubscribe := func() {
		unsubscribeOnce.Do(func() {
			h.mu.Lock()
			delete(h.subscribers, subscriber)
			close(subscriber)
			h.mu.Unlock()
		})
	}

	return subscriber, unsubscribe
}
