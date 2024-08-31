package fifo

import (
	"fmt"
	"time"
)

func (b *Buffer) Read(size int, wait bool) ([]byte, error) {
	if wait {
		for (b.writeIndex-b.readIndex+b.capacity)%b.capacity < size {
			time.Sleep(time.Millisecond)
		}
	} else {
		if (b.writeIndex-b.readIndex+b.capacity)%b.capacity < size {
			return nil, fmt.Errorf("not enough data")
		}
	}

	b.mutex.RLock()
	defer b.mutex.RUnlock()

	packet := make([]byte, size)
	for i := 0; i < size; i++ {
		packet[i] = b.data[(b.readIndex+i)%b.capacity]
	}

	b.readIndex = (b.readIndex + size) % b.capacity
	return packet, nil
}
