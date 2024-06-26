package serial

import (
	"bytes"
	"fmt"
	"io"
	"math"
)

func Filter(port io.ReadWriteCloser, signature []byte) ([]byte, error) {
	header := make([]byte, len(signature))

	for i := 0; i < math.MaxUint8; i++ {
		_, err := port.Read(header)
		if err != nil {
			return nil, err
		}

		if bytes.Equal(header, signature) {
			return header, nil
		}
	}

	return header, fmt.Errorf("failed to filter header")
}
