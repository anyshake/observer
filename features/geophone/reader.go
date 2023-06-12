package geophone

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
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
		return err
	}

	err = binary.Read(
		bytes.NewReader(buffer[:n]),
		binary.LittleEndian,
		options.Geophone,
	)
	if err != nil {
		options.OnErrorCallback(err)
	}

	// Compare checksum
	ok := CompareChecksum(options.Geophone)
	if !ok {
		err = fmt.Errorf("reader: incorrect data frame")
		options.OnErrorCallback(err)
		return nil
	}

	// Low pass filter for vertical
	filteredVertical, err := LowPassFilter(options.Geophone.Vertical[:], FILTER_CUTOFF, FILTER_TAPS)
	if err != nil {
		options.OnErrorCallback(err)
		return err
	}

	// Low pass filter for east-west
	filteredEastWest, err := LowPassFilter(options.Geophone.EastWest[:], FILTER_CUTOFF, FILTER_TAPS)
	if err != nil {
		options.OnErrorCallback(err)
		return err
	}

	// Low pass filter for north-south
	filteredNorthSouth, err := LowPassFilter(options.Geophone.NorthSouth[:], FILTER_CUTOFF, FILTER_TAPS)
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

	// Get synthesis acceleration
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
