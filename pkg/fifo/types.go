package fifo

import "sync"

type Buffer[T any] struct {
	mutex sync.RWMutex

	data       []T
	readIndex  int
	writeIndex int
	capacity   int
}
