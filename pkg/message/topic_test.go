package message

import "testing"

func TestTopicName(t *testing.T) {
	t.Parallel()
	if got := NewBus[int]("events").GetTopicName(); got != "events" {
		t.Fatalf("topic = %q", got)
	}
}
