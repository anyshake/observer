package serial

import (
	"io"
	"log"

	"github.com/albenik/go-serial/v2"
)

func OpenSerial(device string, baud int) io.ReadWriteCloser {
	port, err := serial.Open(device,
		serial.WithBaudrate(baud),
		serial.WithDataBits(8),
		serial.WithParity(serial.NoParity),
		serial.WithStopBits(serial.OneStopBit),
		serial.WithReadTimeout(100),
		serial.WithWriteTimeout(100),
		serial.WithHUPCL(false),
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
