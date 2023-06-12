package serial

import (
	"bytes"
	"fmt"
	"io"
)

func FilterSerial(port io.ReadWriteCloser, signature []byte) error {
	header := make([]byte, len(signature))

	for i := 0; i < 512; i++ {
		port.Read(header)

		if bytes.Equal(header, signature) {
			return nil
		}
	}

	return fmt.Errorf("serial: failed to filter header")
}
