package ringbuf

func (r *Buffer[T]) Values() []T {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	vals := make([]T, r.count)
	for i := 0; i < r.count; i++ {
		vals[i] = r.data[(r.start+i)%r.size]
	}
	return vals
}
