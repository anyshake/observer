package geophone

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"

	"com.geophone.observer/common/serial"
)

func (g *Geophone) Read(port io.ReadWriteCloser, packet *Packet, length int) error {
	// Filter frame header
	err := serial.Filter(port, []byte{0xFC, 0x1B})
	if err != nil {
		return err
	}

	// checksumLength * (uint8 + int32 * length)
	checksumLength := len(packet.Checksum)
	packetSize := checksumLength * (1 + 4*length)

	// Read data frame
	buf := make([]byte, packetSize)
	n, err := serial.Read(port, buf, 5*time.Second)
	if err != nil {
		return err
	}

	// Allocate memory for data frame
	packet.EHZ = make([]int32, length)
	packet.EHE = make([]int32, length)
	packet.EHN = make([]int32, length)

	// Create reader for data frame
	reader := bytes.NewReader(buf[:n])

	// Parse EHZ channel
	err = binary.Read(reader, binary.LittleEndian, packet.EHZ)
	if err != nil {
		return err
	}

	// Parse EHE channel
	err = binary.Read(reader, binary.LittleEndian, packet.EHE)
	if err != nil {
		return err
	}

	// Parse EHN channel
	err = binary.Read(reader, binary.LittleEndian, packet.EHN)
	if err != nil {
		return err
	}

	// Parse checksum
	for i := 0; i < checksumLength; i++ {
		err = binary.Read(reader, binary.LittleEndian, &packet.Checksum[i])
		if err != nil {
			return err
		}
	}

	// Compare checksum
	err = g.isChecksumCorrect(packet)
	if err != nil {
		return err
	}

	return nil
}
