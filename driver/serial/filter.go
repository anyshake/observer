package serial

import (
	"bytes"
	"fmt"
	"io"
)

func Filter(port io.ReadWriteCloser, signature []byte) error {
	header := make([]byte, len(signature))

	for i := 0; i < 64; i++ {
		port.Read(header)

		if bytes.Equal(header, signature) {
			return nil
		}
	}

	return fmt.Errorf("failed to filter header")
}
