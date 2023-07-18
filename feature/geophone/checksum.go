package geophone

import (
	"fmt"
	"unsafe"
)

func (g *Geophone) getChecksum(data []int32) uint8 {
	checksum := uint8(0)

	for i := 0; i < len(data); i++ {
		bytes := (*[4]byte)(unsafe.Pointer(&data[i]))[:]

		for j := 0; j < int(unsafe.Sizeof(int32(0))); j++ {
			checksum ^= bytes[j]
		}
	}

	return checksum
}

func (g *Geophone) isChecksumCorrect(packet *Packet) error {
	var (
		EHZ = g.getChecksum(packet.EHZ[:])
		EHE = g.getChecksum(packet.EHE[:])
		EHN = g.getChecksum(packet.EHN[:])
	)
	if EHZ != packet.Checksum[0] ||
		EHE != packet.Checksum[1] ||
		EHN != packet.Checksum[2] {
		return fmt.Errorf("incorrect packet checksum")
	}

	return nil
}
