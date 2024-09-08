package cache

import "time"

func New(ttl time.Duration) BytesCache {
	return BytesCache{ttl: ttl}
}
