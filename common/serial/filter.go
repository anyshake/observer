package serial

import (
	"bytes"
	"fmt"
	"io"
)

func FilterSerial(port io.ReadWriteCloser, request, signature []byte) error {
	for i := 0; i < 512; i++ {
		header := make([]byte, len(signature))
		port.Write(request)
		port.Read(header)

		if bytes.Equal(header, signature) {
			return nil
		}
	}

	return fmt.Errorf("serial: failed to identify device")
}
