package ringbuf

import "sync"

type Buffer[T any] struct {
	mutex sync.RWMutex

	data  []T
	start int
	count int
	size  int
}
