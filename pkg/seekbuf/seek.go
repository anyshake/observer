package seekbuf

import (
	"fmt"
	"io"
)

func (x *Buffer) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		x.p = int(offset)
	case io.SeekCurrent:
		x.p = x.p + int(offset)
	case io.SeekEnd:
		x.p = x.n + int(offset)
	default:
		return -1, fmt.Errorf("unsupported whence: %d", whence)
	}
	return int64(x.p), nil
}
