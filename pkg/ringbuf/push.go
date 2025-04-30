package ringbuf

func (r *Buffer[T]) Push(val ...T) {
	for _, v := range val {
		r.mutex.Lock()

		if r.count < r.size {
			r.data[(r.start+r.count)%r.size] = v
			r.count++
		} else {
			// overwrite oldest
			r.data[r.start] = v
			r.start = (r.start + 1) % r.size
		}

		r.mutex.Unlock()
	}
}
