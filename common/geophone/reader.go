package geophone

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

func GeophoneReader(port io.ReadWriteCloser, options ReaderOptions) error {
	buffer := make([]byte, 400) // (Length) 100 * (Bytes per uint32_t) 4

	n, err := io.ReadFull(port, buffer)
	if err != nil {
		return err
	}

	err = binary.Read(bytes.NewReader(buffer[:n]), binary.LittleEndian, options.Geophone)
	if err != nil {
		return err
	}

	for i, v := range options.Geophone.Vertical {
		if v > 50000 || v < (-50000) {
			return errors.New("incorrect frame: " + fmt.Sprintf("%d", v))
		}

		options.Acceleration.Vertical[i] = GetAcceleration(v, options.Sensitivity)
		options.Acceleration.Timestamp = GetTimestamp()
	}

	return nil
}

func ReaderDaemon(device string, baud int, options ReaderOptions) {
	port := OpenGeophone(device, baud)
	defer port.Close()

	for {
		err := GeophoneReader(port, options)
		if err != nil {
			options.OnErrorCallback(err)
			CloseGeophone(port)
			port = OpenGeophone(device, baud)

			continue
		}

		options.OnDataCallback(options.Acceleration)
	}
}
