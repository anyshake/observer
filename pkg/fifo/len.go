package fifo

func (b *Buffer[T]) Len() int {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return (b.writeIndex - b.readIndex + b.capacity) % b.capacity
}
