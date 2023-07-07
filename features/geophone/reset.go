package geophone

import (
	"io"

	"com.geophone.observer/common/serial"
)

func ResetGeophone(port io.ReadWriteCloser) error {
	_, err := port.Write([]byte{0x61})
	if err != nil {
		return err
	}

	err = serial.FilterSerial(port, []byte{0xAC, 0x55})
	if err != nil {
		return err
	}

	return nil
}
