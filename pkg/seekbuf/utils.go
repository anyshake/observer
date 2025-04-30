package seekbuf

func (x *Buffer) grow(n int) {
	if n >= x.b.Cap() {
		b := x.Bytes()
		x.b.Grow(n)
		copy(x.Bytes()[:x.n], b)
	}
}

func (x *Buffer) Reset() {
	x.n = 0
	x.p = 0
}

func (x *Buffer) String() string {
	return string(x.Bytes())
}

func (x *Buffer) Bytes() []byte {
	return x.b.Bytes()[:x.n]
}

func (x *Buffer) Len() int {
	return x.n
}
