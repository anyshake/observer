package geophone

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"
	"unsafe"

	"com.geophone.observer/common/serial"
)

func (g *Geophone) Read(port io.ReadWriteCloser, packet *Packet) error {
	// Filter frame header
	err := serial.Filter(port, []byte{0xFC, 0x1B})
	if err != nil {
		return err
	}

	// Read data frame
	buf := make([]byte, unsafe.Sizeof(Packet{}))
	n, err := serial.Read(port, buf, 5*time.Second)
	if err != nil {
		return err
	}

	// Parse data frame
	err = binary.Read(
		bytes.NewReader(buf[:n]),
		binary.LittleEndian, packet,
	)
	if err != nil {
		return err
	}

	// Compare checksum
	err = g.isChecksumCorrect(packet)
	if err != nil {
		return err
	}

	return nil
}
