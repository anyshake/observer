package geophone

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
	"time"
	"unsafe"

	"com.geophone.observer/common/serial"
)

func GeophoneReader(port io.ReadWriteCloser, options GeophoneOptions) error {
	buffer := make([]byte, unsafe.Sizeof(Geophone{}))

	err := serial.FilterSerial(port, []byte{0x55, 0x55}, []byte{0x55, 0xAA})
	if err != nil {
		return err
	}

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

	v := reflect.ValueOf(options.Geophone).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if field.Type.Kind() == reflect.Float32 {
			if math.Abs(v.Field(i).Float()) > 1000 {
				err = fmt.Errorf("reader: incorrect data frame")
				options.OnErrorCallback(err)
				return err
			}
		}
	}

	options.Acceleration = &Acceleration{
		Vertical: GetAcceleration(
			float64(options.Geophone.Vertical),
			options.Sensitivity.Vertical,
		),
		EastWest: GetAcceleration(
			float64(options.Geophone.EastWest),
			options.Sensitivity.EastWest,
		),
		NorthSouth: GetAcceleration(
			float64(options.Geophone.NorthSouth),
			options.Sensitivity.NorthSouth,
		),
	}

	options.Acceleration.Synthesis = GetSynthesis(
		options.Acceleration.Vertical,
		options.Acceleration.EastWest,
		options.Acceleration.NorthSouth,
	)

	options.OnDataCallback(options.Acceleration)
	return nil
}

func ReaderDaemon(device string, baud int, options GeophoneOptions) {
	port := serial.OpenSerial(device, baud)
	defer serial.CloseSerial(port)

	for {
		err := GeophoneReader(port, options)
		if err != nil {
			serial.CloseSerial(port)
			options.OnErrorCallback(err)
			time.Sleep(100 * time.Millisecond)
			port = serial.OpenSerial(device, baud)

			continue
		}

		time.Sleep(15 * time.Millisecond)
	}
}
