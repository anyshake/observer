package serial

import "io"

func Close(port io.ReadWriteCloser) error {
	if port == nil {
		return nil
	}

	return port.Close()
}
