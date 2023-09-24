package serial

import (
	"bytes"
	"fmt"
	"io"
)

func Filter(port io.ReadWriteCloser, signature []byte, retry int) ([]byte, error) {
	header := make([]byte, len(signature))

	for i := 0; i < retry; i++ {
		port.Read(header)

		if bytes.Equal(header, signature) {
			return nil, nil
		}
	}

	return header, fmt.Errorf("failed to filter header")
}
