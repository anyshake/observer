package ringbuf

func (r *Buffer[T]) Push(val ...T) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	n := len(val)
	if n == 0 || r.size <= 0 {
		return
	}

	// A longer batch replaces the buffer with its own tail, in order.
	if n >= r.size {
		copy(r.data, val[n-r.size:])
		r.start = 0
		r.count = r.size
		return
	}

	end := (r.start + r.count) % r.size
	space := r.size - end
	if n <= space {
		copy(r.data[end:end+n], val)
	} else {
		copy(r.data[end:], val[:space])
		copy(r.data, val[space:])
	}

	if r.count+n <= r.size {
		r.count += n
		return
	}

	overflow := r.count + n - r.size
	r.start = (r.start + overflow) % r.size
	r.count = r.size
}
