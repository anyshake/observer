package seekbuf

import "errors"

func (x *Buffer) Write(p []byte) (int, error) {
	n := len(p)
	if n == 0 {
		return 0, nil
	}
	if n > int(^uint(0)>>1)-x.p {
		return 0, errors.New("buffer size overflow")
	}
	t := x.p + n
	x.grow(t)
	if x.p > x.n {
		clear(x.b.Bytes()[x.n:x.p])
	}
	copy(x.Bytes()[x.p:t], p)
	if t > x.n {
		x.n = t
	}
	x.p = t
	return n, nil
}

func (x *Buffer) WriteString(s string) (int, error) {
	return x.Write([]byte(s))
}
