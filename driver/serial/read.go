package serial

import (
	"fmt"
	"io"
	"time"
)

func Read(r io.Reader, buf []byte, timeout time.Duration) (n int, err error) {
	min := len(buf)

	if len(buf) < min {
		return 0, io.ErrShortBuffer
	}

	start := time.Now()
	for n < min {
		if time.Since(start) > timeout {
			return 0, fmt.Errorf("timeout due to no response")
		}

		nn, err := r.Read(buf[n:])
		if err != nil {
			return 0, err
		}

		n += nn
	}

	if n >= min {
		err = nil
	} else if n > 0 && err == io.EOF {
		err = io.ErrUnexpectedEOF
	}

	return n, err
}
