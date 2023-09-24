package geophone

import (
	"fmt"
	"io"

	"github.com/bclswl0827/observer/driver/serial"
)

func (g *Geophone) Reset(port io.ReadWriteCloser) error {
	_, err := port.Write(RESET_WORD[:])
	if err != nil {
		return err
	}

	_, err = serial.Filter(port, ACK_WORD[:], 64)
	if err != nil {
		return fmt.Errorf("failed to reset geophone")
	}
	return nil
}
