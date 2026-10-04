package seekbuf

import "io"

func (x *Buffer) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if x.p >= x.n {
		return 0, io.EOF
	}
	n := copy(p, x.Bytes()[x.p:])
	x.p += n
	return n, nil
}
