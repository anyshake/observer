package fifo

func (b *Buffer[T]) Reset() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.readIndex = 0
	b.writeIndex = 0
}
