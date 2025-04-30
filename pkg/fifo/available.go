package fifo

func (b *Buffer[T]) Available(size int) bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return (b.writeIndex-b.readIndex+b.capacity)%b.capacity > size
}
