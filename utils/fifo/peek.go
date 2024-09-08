package fifo

import (
	"fmt"
)

func (b *Buffer) Peek(header []byte, size int) ([]byte, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	for {
		if (b.writeIndex-b.readIndex+b.capacity)%b.capacity < size {
			return nil, fmt.Errorf("not enough data")
		}

		isHeaderFind := true
		for i := 0; i < len(header); i++ {
			if b.data[(b.readIndex+i)%b.capacity] != header[i] {
				isHeaderFind = false
				break
			}
		}

		if isHeaderFind {
			break
		}

		b.readIndex = (b.readIndex + 1) % b.capacity
	}

	packet := make([]byte, size)
	for i := 0; i < size; i++ {
		packet[i] = b.data[(b.readIndex+i)%b.capacity]
	}

	b.readIndex = (b.readIndex + size) % b.capacity
	return packet, nil
}
