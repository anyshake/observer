package fifo

func New[T any](size int) Buffer[T] {
	return Buffer[T]{
		data:     make([]T, size),
		capacity: size,
	}
}
