package serial

import (
	"bytes"
	"fmt"
	"io"
)

func FilterSerial(port io.ReadWriteCloser, padding, signature []byte) error {
	header := make([]byte, len(signature))

	for i := 0; i < len(padding); i++ {
		port.Read(header)

		if bytes.Equal(header, signature) {
			return nil
		}

		if bytes.Equal(header, padding) || bytes.Contains(header, padding[:1]) {
			continue
		}

	}

	return fmt.Errorf("serial: failed to filter header")
}
