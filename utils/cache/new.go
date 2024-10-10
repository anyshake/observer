package cache

import "time"

func New(ttl time.Duration) AnyCache {
	return AnyCache{ttl: ttl}
}
