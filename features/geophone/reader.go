package geophone

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
	"time"
	"unsafe"

	"com.geophone.observer/common/serial"
)

func GeophoneReader(port io.ReadWriteCloser, options GeophoneOptions) error {
	serial.FilterSerial(port, []byte{0x55, 0xAA})

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

	val := reflect.ValueOf(options.Geophone).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.Slice {
			for j := 0; j < field.Len(); j++ {
				if field.Index(j).Int() > 50000 || field.Index(j).Int() < (-50000) {
					options.OnErrorCallback(
						errors.New("incorrect frame: " + fmt.Sprintf("%d", field.Index(j).Int())),
					)
					return errors.New("incorrect frame: " + fmt.Sprintf("%d", field.Index(j).Int()))
				}
			}
		}
	}

	for i, v := range options.Geophone.Vertical {
		options.Acceleration.Vertical[i] = GetAcceleration(v, options.Sensitivity.Vertical)
	}

	// for i, v := range options.Geophone.EastWest {
	// 	options.Acceleration.EastWest[i] = GetAcceleration(v, options.Sensitivity.EastWest)
	// }

	// for i, v := range options.Geophone.NorthSouth {
	// 	options.Acceleration.NorthSouth[i] = GetAcceleration(v, options.Sensitivity.NorthSouth)
	// }

	if options.Geophone.Latitude == -1 &&
		options.Geophone.Longitude == -1 {
		options.Acceleration = &Acceleration{
			Altitude:   options.LocationFallback.Altitude,
			Latitude:   options.LocationFallback.Latitude,
			Longitude:  options.LocationFallback.Longitude,
			Vertical:   options.Acceleration.Vertical,
			EastWest:   options.Acceleration.EastWest,
			NorthSouth: options.Acceleration.NorthSouth,
			Synthesis: GetSynthesis(
				options.Acceleration.Vertical,
				options.Acceleration.EastWest,
				options.Acceleration.NorthSouth,
			),
		}
	} else {
		options.Acceleration = &Acceleration{
			Altitude:   float64(options.Geophone.Altitude),
			Latitude:   float64(options.Geophone.Latitude),
			Longitude:  float64(options.Geophone.Longitude),
			Vertical:   options.Acceleration.Vertical,
			EastWest:   options.Acceleration.EastWest,
			NorthSouth: options.Acceleration.NorthSouth,
			Synthesis: GetSynthesis(
				options.Acceleration.Vertical,
				options.Acceleration.EastWest,
				options.Acceleration.NorthSouth,
			),
		}
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
			time.Sleep(time.Second)
			port = serial.OpenSerial(device, baud)

			continue
		}
	}
}
