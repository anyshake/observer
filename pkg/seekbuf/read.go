package seekbuf

func (x *Buffer) Read(p []byte) (int, error) {
	return copy(p, x.Bytes()[x.p:]), nil
}
