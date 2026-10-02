package ringbuf

func (r *Buffer[T]) Values() []T {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	vals := make([]T, r.count)
	for i := 0; i < r.count; i++ {
		vals[i] = r.data[(r.start+i)%r.size]
	}
	return vals
}
