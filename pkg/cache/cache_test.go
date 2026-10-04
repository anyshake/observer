package cache_test

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/anyshake/observer/pkg/cache"
)

func TestGenericCacheLifecycle(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := cache.NewGeneric[string](time.Minute)
		if c.Valid() || c.Get() != "" {
			t.Fatal("new cache is not empty and invalid")
		}
		c.Set("first")
		if !c.Valid() || c.Get() != "first" {
			t.Fatal("Set() did not populate cache")
		}
		time.Sleep(time.Minute - time.Nanosecond)
		if !c.Valid() {
			t.Fatal("cache expired before TTL")
		}
		time.Sleep(time.Nanosecond)
		if c.Valid() {
			t.Fatal("cache remains valid at TTL")
		}
		c.Set("second")
		if !c.Valid() || c.Get() != "second" {
			t.Fatal("Set() did not refresh expired cache")
		}
		c.Clear()
		if c.Valid() || c.Get() != "" {
			t.Fatal("Clear() did not invalidate and empty cache")
		}
		c.Set("")
		if !c.Valid() {
			t.Fatal("cached zero value should still be valid")
		}
	})
}

func TestKvCacheLifecycle(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := cache.NewKv[int](time.Minute)
		if c.Valid() {
			t.Fatal("new cache is valid")
		}
		if got, found := c.Get("missing"); found || got != 0 {
			t.Fatalf("missing key = %d, %v", got, found)
		}
		c.Set("first", 1)
		time.Sleep(30 * time.Second)
		c.Set("second", 0)
		time.Sleep(30 * time.Second)
		if !c.Valid() {
			t.Fatal("Set() did not refresh the shared cache TTL")
		}
		if got, found := c.Get("first"); !found || got != 1 {
			t.Fatalf("first key = %d, %v", got, found)
		}
		if got, found := c.Get("second"); !found || got != 0 {
			t.Fatalf("zero value key = %d, %v", got, found)
		}
		time.Sleep(30 * time.Second)
		if c.Valid() {
			t.Fatal("cache remains valid at TTL")
		}
		c.Clear()
		if _, found := c.Get("first"); found || c.Valid() {
			t.Fatal("Clear() retained cache entries or validity")
		}
		c.Set("first", 2)
		if got, found := c.Get("first"); !found || got != 2 || !c.Valid() {
			t.Fatal("cache cannot be reused after Clear()")
		}
	})
}

func TestCacheTTLBoundaries(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		for _, ttl := range []time.Duration{0, -time.Second, time.Duration(1<<63 - 1)} {
			generic := cache.NewGeneric[int](ttl)
			kv := cache.NewKv[int](ttl)
			if generic.Valid() || kv.Valid() {
				t.Fatalf("empty cache valid with TTL %v", ttl)
			}
			generic.Set(1)
			kv.Set("key", 1)
			if generic.Valid() != (ttl > 0) || kv.Valid() != (ttl > 0) {
				t.Fatalf("unexpected validity with TTL %v", ttl)
			}
			generic.Clear()
			kv.Clear()
			if generic.Valid() || kv.Valid() {
				t.Fatalf("cleared cache valid with TTL %v", ttl)
			}
		}
	})
}
