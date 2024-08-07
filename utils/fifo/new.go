package fifo

func New(size int) Buffer {
	return Buffer{
		data:     make([]byte, size),
		capacity: size,
	}
}
