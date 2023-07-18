package geophone

import "time"

const (
	// Maximum value 375
	PACKET_SIZE int = 10
	// READY_THRESHOLD should be smaller than 1 second
	READY_THRESHOLD time.Duration = 900 * time.Millisecond
	// TIMEOUT_THRESHOLD should be greater than 1 second
	TIMEOUT_THRESHOLD time.Duration = 3 * time.Second
)

const MODULE string = "geophone"

type Geophone struct{}

type Packet struct {
	EHZ      [PACKET_SIZE]int32 // Vertical
	EHE      [PACKET_SIZE]int32 // East-West
	EHN      [PACKET_SIZE]int32 // North-South
	Checksum [3]uint8
}
