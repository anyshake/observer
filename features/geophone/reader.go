package geophone

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"unsafe"
)

func GeophoneReader(port io.ReadWriteCloser, options GeophoneOptions) error {
	for {
		header := make([]byte, 2)
		_, err := port.Read(header)
		if err != nil {
			options.OnErrorCallback(err)
			return err
		}

		if header[0] == 0x55 &&
			header[1] == 0xAA {
			break
		}
	}

	buffer := make([]byte, unsafe.Sizeof(Geophone{}))
	n, err := io.ReadFull(port, buffer)
	if err != nil {
		options.OnErrorCallback(err)
		return err
	}

	err = binary.Read(
		bytes.NewReader(buffer[:n]),
		binary.LittleEndian,
		options.Geophone,
	)
	if err != nil {
		options.OnErrorCallback(err)
		return err
	}

	for i, v := range options.Geophone.Vertical {
		if v > 50000 || v < (-50000) {
			options.OnErrorCallback(
				errors.New("incorrect frame: " + fmt.Sprintf("%d", v)),
			)
			return errors.New("incorrect frame: " + fmt.Sprintf("%d", v))
		}

		options.Acceleration.Vertical[i] = GetAcceleration(v, options.Sensitivity)
	}

	if options.Geophone.Latitude == -1 &&
		options.Geophone.Longitude == -1 {
		options.Acceleration = &Acceleration{
			Altitude:  options.LocationFallback.Altitude,
			Latitude:  options.LocationFallback.Latitude,
			Longitude: options.LocationFallback.Longitude,
			Vertical:  options.Acceleration.Vertical,
		}
	} else {
		options.Acceleration = &Acceleration{
			Altitude:  float64(options.Geophone.Altitude),
			Latitude:  float64(options.Geophone.Latitude),
			Longitude: float64(options.Geophone.Longitude),
			Vertical:  options.Acceleration.Vertical,
		}
	}
	options.OnDataCallback(options.Acceleration)

	return nil
}

func ReaderDaemon(device string, baud int, options GeophoneOptions) {
	port := OpenGeophone(device, baud)
	defer CloseGeophone(port)

	for {
		err := GeophoneReader(port, options)
		if err != nil {
			CloseGeophone(port)
			port = OpenGeophone(device, baud)

			continue
		}
	}
}
