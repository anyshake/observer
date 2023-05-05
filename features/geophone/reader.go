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
	err := serial.FilterSerial(port,
		[]byte{0xAA, 0x55},
	)
	if err != nil {
		return err
	}

	buffer := make([]byte, unsafe.Sizeof(Geophone{}))
	n, err := serial.ReadSerial(port, buffer, 2*time.Second)
	if err != nil {
		options.OnErrorCallback(err)
	}

	err = binary.Read(
		bytes.NewReader(buffer[:n]),
		binary.LittleEndian,
		options.Geophone,
	)
	if err != nil {
		options.OnErrorCallback(err)
	}

	val := reflect.ValueOf(options.Geophone).Elem()
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		if fieldVal.Kind() == reflect.Array {
			for j := 0; j < fieldVal.Len(); j++ {
				itemVal := fieldVal.Index(j)

				if math.Abs(itemVal.Float()) > 10000 {
					err = fmt.Errorf("reader: incorrect data frame")
					options.OnErrorCallback(err)
					return nil
				}

				if err == nil {
					switch val.Type().Field(i).Name {
					case "Vertical":
						options.Acceleration.Vertical[j] = GetAcceleration(itemVal.Float(), options.Sensitivity.Vertical)
					case "EastWest":
						options.Acceleration.EastWest[j] = GetAcceleration(itemVal.Float(), options.Sensitivity.EastWest)
					case "NorthSouth":
						options.Acceleration.NorthSouth[j] = GetAcceleration(itemVal.Float(), options.Sensitivity.NorthSouth)
					}
				}
			}
		}
	}

	for i := range options.Acceleration.Vertical {
		options.Acceleration.Synthesis[i] = GetSynthesis(
			options.Acceleration.Vertical[i],
			options.Acceleration.EastWest[i],
			options.Acceleration.NorthSouth[i],
		)
	}

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
			time.Sleep(500 * time.Millisecond)
			port = serial.OpenSerial(device, baud)

			continue
		}
	}
}
