package seekbuf

func (x *Buffer) Write(p []byte) (int, error) {
	n := len(p)
	t := x.p + n
	x.grow(t)
	copy(x.Bytes()[x.p:t], p)
	if t > x.n {
		x.n = t
	}
	x.p = t
	return n, nil
}

func (x *Buffer) WriteString(s string) (int, error) {
	n := len(s)
	t := x.p + n
	x.grow(t)
	copy(x.Bytes()[x.p:t], s)
	if t > x.n {
		x.n = t
	}
	x.p = t
	return n, nil
}
