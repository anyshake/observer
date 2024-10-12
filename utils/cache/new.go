package cache

import "time"

func New(ttl time.Duration) AnyCache {
	return AnyCache{ttl: ttl}
}

func NewKv(ttl time.Duration) KvCache {
	return KvCache{
		cache: map[any]any{},
		ttl:   ttl,
	}
}
