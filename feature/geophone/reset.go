package geophone

import (
	"fmt"
	"io"

	"github.com/bclswl0827/observer/driver/serial"
)

func (g *Geophone) Reset(port io.ReadWriteCloser) error {
	_, err := port.Write([]byte{0x61})
	if err != nil {
		return err
	}

	err = serial.Filter(port, []byte{0xFC, 0x2B})
	if err != nil {
		return fmt.Errorf("failed to reset geophone")
	}
	return nil
}
