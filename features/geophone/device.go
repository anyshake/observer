package geophone

import (
	"io"
	"log"

	"github.com/tarm/serial"
)

func OpenGeophone(device string, baud int) io.ReadWriteCloser {
	options := &serial.Config{
		Name: device,
		Baud: baud,
	}

	port, err := serial.OpenPort(options)
	if err != nil {
		log.Fatalln(err)
	}

	return port
}

func CloseGeophone(port io.ReadWriteCloser) error {
	err := port.Close()
	if err != nil {
		return err
	}

	return nil
}
