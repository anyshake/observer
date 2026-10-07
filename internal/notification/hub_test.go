package notification

import "testing"

func TestHubPublishAndUnsubscribe(t *testing.T) {
	t.Parallel()

	hub := NewHub(0)
	events, unsubscribe := hub.Subscribe()
	hub.Publish(Event{ServiceID: "quakesense", Message: "triggered", Level: LevelWarning})
	event := <-events
	if event.ID == "" || event.OccurredAt == 0 || event.Message != "triggered" {
		t.Fatalf("event = %+v", event)
	}

	unsubscribe()
	unsubscribe()
	hub.Publish(Event{ID: "fixed", OccurredAt: 10, Message: "after"})
	if _, ok := <-events; ok {
		t.Fatal("subscriber stayed open after unsubscribe")
	}

	full, unsubscribeFull := hub.Subscribe()
	hub.Publish(Event{ID: "one", OccurredAt: 1})
	hub.Publish(Event{ID: "dropped", OccurredAt: 2})
	if got := <-full; got.ID != "one" {
		t.Fatalf("buffered event = %+v", got)
	}
	unsubscribeFull()
}
