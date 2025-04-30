package fifo

import (
	"errors"
)

func (b *Buffer[T]) Peek(header []T, size int) ([]T, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	for {
		if (b.writeIndex-b.readIndex+b.capacity)%b.capacity < size {
			return nil, errors.New("not enough data")
		}

		isHeaderFind := true
		if len(header) > 0 {
			for i := 0; i < len(header); i++ {
				if any(b.data[(b.readIndex+i)%b.capacity]) != any(header[i]) {
					isHeaderFind = false
					break
				}
			}
		}

		if isHeaderFind {
			break
		}

		b.readIndex = (b.readIndex + 1) % b.capacity
	}

	packet := make([]T, size)
	for i := 0; i < size; i++ {
		packet[i] = b.data[(b.readIndex+i)%b.capacity]
	}

	b.readIndex = (b.readIndex + size) % b.capacity
	return packet, nil
}
