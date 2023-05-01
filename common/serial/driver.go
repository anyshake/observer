package serial

import (
	"io"
	"log"

	"github.com/albenik/go-serial/v2"
)

func OpenSerial(device string, baud int) io.ReadWriteCloser {
	port, err := serial.Open(device,
		serial.WithHUPCL(false),
		serial.WithDataBits(8),
		serial.WithBaudrate(baud),
		serial.WithReadTimeout(500),
		serial.WithWriteTimeout(500),
		serial.WithParity(serial.NoParity),
		serial.WithStopBits(serial.OneStopBit),
	)
	if err != nil {
		log.Fatalln(err)
	}

	port.SetDTR(true)
	port.SetRTS(true)

	return port
}

func CloseSerial(port io.ReadWriteCloser) error {
	return port.Close()
}
