package ringbuf

func (r *Buffer[T]) Push(val ...T) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, v := range val {
		if r.count < r.size {
			r.data[(r.start+r.count)%r.size] = v
			r.count++
		} else {
			// overwrite oldest
			r.data[r.start] = v
			r.start = (r.start + 1) % r.size
		}
	}
}
