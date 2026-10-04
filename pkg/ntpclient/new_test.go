package ntpclient

import (
	"slices"
	"testing"
	"time"
)

func TestNewPreservesServerPool(t *testing.T) {
	t.Parallel()
	now := time.Unix(1700000000, 0)
	client, err := New([]string{"ntp://time.example.test:123", "ntp://[2001:db8::1]:124", "ntp://pool.example.test"}, 2, 5, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"time.example.test:123", "[2001:db8::1]:124", "pool.example.test"}
	if !slices.Equal(client.pool, want) {
		t.Fatalf("server pool = %v, want %v", client.pool, want)
	}
	if client.retries != 2 || client.readTimeout != 5*time.Second || !client.timeFunc().Equal(now) {
		t.Fatal("client options were not preserved")
	}
}

func TestNewRejectsInvalidPool(t *testing.T) {
	t.Parallel()
	for _, pool := range [][]string{nil, {""}, {"://bad"}, {"ntp://"}, {"ntp:///path"}, {"https://time.example.test"}, {"time.example.test"}, {"ntp://host:bad"}, {"ntp://host:0"}, {"ntp://host:65536"}} {
		if client, err := New(pool, 1, 1, nil); err == nil || client != nil {
			t.Errorf("New(%q) = %v, %v, want error", pool, client, err)
		}
	}
}

func TestNewDeduplicatesEndpoints(t *testing.T) {
	t.Parallel()
	client, err := New([]string{
		"ntp://Time.Example.Test.",
		"ntp://time.example.test:123",
		"ntp://time.example.test:00123",
		"ntp://time.example.test:124",
		"ntp://[2001:db8::1]",
		"ntp://[2001:0db8:0:0:0:0:0:1]:123",
		"ntp://192.0.2.1",
		"ntp://[::ffff:192.0.2.1]:123",
	}, 0, 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Time.Example.Test.", "time.example.test:124", "[2001:db8::1]", "192.0.2.1"}
	if !slices.Equal(client.pool, want) {
		t.Fatalf("deduplicated pool = %v, want %v", client.pool, want)
	}
	if len(client.servers) != len(want) {
		t.Fatalf("poll states = %d, want %d", len(client.servers), len(want))
	}
	for _, server := range want {
		if client.servers[server].interval != minPollInterval {
			t.Errorf("%s: initial poll interval = %v", server, client.servers[server].interval)
		}
	}
}

func TestNewRejectsNegativeOptions(t *testing.T) {
	t.Parallel()
	for _, options := range [][2]int{{-1, 1}, {1, -1}} {
		if client, err := New([]string{"ntp://time.example.test"}, options[0], options[1], nil); err == nil || client != nil {
			t.Errorf("New with retries=%d, timeout=%d = %v, %v, want error", options[0], options[1], client, err)
		}
	}
}

func TestNewDefaultsTimeFunction(t *testing.T) {
	t.Parallel()
	client, err := New([]string{"ntp://time.example.test"}, 1, 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	before := time.Now()
	got := client.timeFunc()
	if got.Before(before) || got.After(time.Now()) {
		t.Fatalf("default clock = %v, outside current time range", got)
	}
}

func TestQueryAverageRejectsNonPositiveAttempts(t *testing.T) {
	t.Parallel()
	client := &Client{}
	for _, attempts := range []int{0, -1} {
		if _, err := client.QueryAverage(attempts); err == nil {
			t.Errorf("QueryAverage(%d) accepted invalid attempts", attempts)
		}
	}
}
