package serial

import (
	"io"
	"log"

	"github.com/tarm/serial"
)

func OpenSerial(device string, baud int) io.ReadWriteCloser {
	options := &serial.Config{
		Name: device,
		Baud: baud,
	}

	port, err := serial.OpenPort(options)
	if err != nil {
		log.Fatalln(err)
	}

	port.Flush()
	return port
}

func CloseSerial(port io.ReadWriteCloser) error {
	return port.Close()
}
