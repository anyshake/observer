package serial

import (
	"bytes"
	"io"
)

func FilterSerial(port io.ReadWriteCloser, request, signature []byte) {
	for {
		header := make([]byte, len(signature))
		port.Write(request)
		port.Read(header)

		if bytes.Equal(header, signature) {
			break
		}
	}
}
