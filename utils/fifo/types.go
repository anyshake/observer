package fifo

import "sync"

type Buffer struct {
	data       []byte
	readIndex  int
	writeIndex int
	capacity   int
	mutex      sync.Mutex
}
