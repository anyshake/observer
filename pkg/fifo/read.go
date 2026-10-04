package fifo

import "errors"

func (b *Buffer[T]) Read(size int) ([]T, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if size < 0 {
		return nil, errors.New("invalid read size")
	}

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
