package fifo

import "errors"

func (b *Buffer[T]) Read(size int) ([]T, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if (b.writeIndex-b.readIndex+b.capacity)%b.capacity < size {
		return nil, errors.New("not enough data")
	}

	packet := make([]T, size)
	for i := 0; i < size; i++ {
		packet[i] = b.data[(b.readIndex+i)%b.capacity]
	}

	b.readIndex = (b.readIndex + size) % b.capacity
	return packet, nil
}
