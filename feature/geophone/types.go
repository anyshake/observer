package geophone

import "time"

const (
	// READY_THRESHOLD should be smaller than 1 second
	READY_THRESHOLD time.Duration = 900 * time.Millisecond
	// TIMEOUT_THRESHOLD should be greater than 1 second
	TIMEOUT_THRESHOLD time.Duration = 3 * time.Second
)

const MODULE string = "geophone"

type Geophone struct{}

type Packet struct {
	EHZ      []int32 // Vertical
	EHE      []int32 // East-West
	EHN      []int32 // North-South
	Checksum [3]uint8
}
