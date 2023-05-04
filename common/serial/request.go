package serial

import (
	"io"
)

func RequestSerial(port io.ReadWriteCloser, request []byte) error {
	_, err := port.Write(request)
	if err != nil {
		return err
	}

	return nil
}
