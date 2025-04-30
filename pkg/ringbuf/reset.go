package ringbuf

func (r *Buffer[T]) Reset() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.start = 0
	r.count = 0
}
