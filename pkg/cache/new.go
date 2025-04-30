package cache

import "time"

func New(ttl time.Duration) AnyCache {
	return AnyCache{ttl: ttl}
}

func NewKv[T any](ttl time.Duration) KvCache[T] {
	return KvCache[T]{
		cache: map[any]T{},
		ttl:   ttl,
	}
}
