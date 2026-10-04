package seekbuf

import (
	"fmt"
	"io"
)

func (x *Buffer) Seek(offset int64, whence int) (int64, error) {
	var base int64
	switch whence {
	case io.SeekStart:
	case io.SeekCurrent:
		base = int64(x.p)
	case io.SeekEnd:
		base = int64(x.n)
	default:
		return -1, fmt.Errorf("unsupported whence: %d", whence)
	}
	if offset < -base || offset > int64(^uint(0)>>1)-base {
		return -1, fmt.Errorf("invalid seek offset: %d", offset)
	}
	x.p = int(base + offset)
	return int64(x.p), nil
}
