package serial

import (
	"io"

	"github.com/bclswl0827/go-serial"
)

func Open(device string, baud int) (io.ReadWriteCloser, error) {
	port, err := serial.Open(device,
		serial.WithHUPCL(true),
		serial.WithDataBits(8),
		serial.WithBaudrate(baud),
		serial.WithReadTimeout(5),
		serial.WithWriteTimeout(5),
		serial.WithParity(serial.NoParity),
		serial.WithStopBits(serial.OneStopBit),
	)
	if err != nil {
		return nil, err
	}

	port.SetDTR(true)
	port.SetRTS(true)

	return port, nil
}
