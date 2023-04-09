package serial

import (
	"bytes"
	"io"
)

func FilterSerial(port io.ReadWriteCloser, signature []byte) {
	for {
		header := make([]byte, len(signature))
		port.Read(header)

		if bytes.Equal(header, signature) {
			break
		}
	}
}
