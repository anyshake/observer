package ringbuf

func New[T any](size int) *Buffer[T] {
	return &Buffer[T]{
		data: make([]T, size),
		size: size,
	}
}
