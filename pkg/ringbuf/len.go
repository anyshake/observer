package ringbuf

func (r *Buffer[T]) Len() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.count
}
