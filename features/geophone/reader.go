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

	var (
		rawVertical   = make([]float64, PACKET_SIZE)
		rawEastWest   = make([]float64, PACKET_SIZE)
		rawNorthSouth = make([]float64, PACKET_SIZE)
	)

	// Error checking
	val := reflect.ValueOf(options.Geophone).Elem()
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)

		if fieldVal.Kind() == reflect.Array {
			for j := 0; j < fieldVal.Len(); j++ {
				itemVal := fieldVal.Index(j)

				if math.Abs(itemVal.Float()) > 100 {
					err = fmt.Errorf("reader: incorrect data frame")
					options.OnErrorCallback(err)
					return nil
				}

				switch val.Type().Field(i).Name {
				case "Vertical":
					rawVertical[j] = itemVal.Float()
				case "EastWest":
					rawEastWest[j] = itemVal.Float()
				case "NorthSouth":
					rawNorthSouth[j] = itemVal.Float()
				}
			}
		}
	}

	// Low pass filter for vertical
	filteredVertical, err := LowPassFilter(rawVertical, FILTER_CUTOFF, FILTER_TAPS)
	if err != nil {
		options.OnErrorCallback(err)
		return err
	}

	// Low pass filter for east-west
	filteredEastWest, err := LowPassFilter(rawEastWest, FILTER_CUTOFF, FILTER_TAPS)
	if err != nil {
		options.OnErrorCallback(err)
		return err
	}

	// Low pass filter for north-south
	filteredNorthSouth, err := LowPassFilter(rawNorthSouth, FILTER_CUTOFF, FILTER_TAPS)
	if err != nil {
		options.OnErrorCallback(err)
		return err
	}

	// Get vertical acceleration
	for i := 0; i < PACKET_SIZE; i++ {
		options.Acceleration.Vertical[i] = GetAcceleration(
			filteredVertical[i], options.Sensitivity.Vertical,
		)
	}

	// Get east-west acceleration
	for i := 0; i < PACKET_SIZE; i++ {
		options.Acceleration.EastWest[i] = GetAcceleration(
			filteredEastWest[i], options.Sensitivity.EastWest,
		)
	}

	// Get north-south acceleration
	for i := 0; i < PACKET_SIZE; i++ {
		options.Acceleration.NorthSouth[i] = GetAcceleration(
			filteredNorthSouth[i], options.Sensitivity.NorthSouth,
		)
	}

	for i := 0; i < PACKET_SIZE; i++ {
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
