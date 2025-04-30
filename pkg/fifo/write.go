package fifo

func (b *Buffer[T]) Write(p ...T) (n int, err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for _, bt := range p {
		b.data[b.writeIndex] = bt

		b.writeIndex = (b.writeIndex + 1) % b.capacity
		if b.writeIndex == b.readIndex {
			b.readIndex = (b.readIndex + 1) % b.capacity
		}

		n++
	}

	return n, nil
}
