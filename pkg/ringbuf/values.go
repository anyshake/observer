package ringbuf

func (r *Buffer[T]) Values() []T {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	vals := make([]T, r.count)
	if r.count == 0 {
		return vals
	}

	tail := r.size - r.start
	if r.count <= tail {
		copy(vals, r.data[r.start:r.start+r.count])
		return vals
	}
	copy(vals, r.data[r.start:])
	copy(vals[tail:], r.data[:r.count-tail])
	return vals
}
