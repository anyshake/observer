package geophone

import "time"

const MODULE string = "geophone"

const (
	// READY_THRESHOLD should be smaller than 1 second
	READY_THRESHOLD time.Duration = 900 * time.Millisecond
	// TIMEOUT_THRESHOLD should be greater than 1 second
	TIMEOUT_THRESHOLD time.Duration = 3 * time.Second
)

const (
	// TARGET_DAMPING is the target damping ratio
	TARGET_DAMPING float64 = 0.707
	// TARGET_FREQUENCY is the target frequency
	TARGET_FREQUENCY float64 = 0.1
)

var (
	// RESET_WORD resets geophone ADC module
	RESET_WORD = [...]byte{0x61}
	// SYNC_WORD indicates a data packet is following
	SYNC_WORD = [...]byte{0xFC, 0x1B}
	// ACK_WORD indicates a valid command is received
	ACK_WORD = [...]byte{0xFC, 0x2B}
)

type Geophone struct{}

type Packet struct {
	EHZ      []int32 // Vertical
	EHE      []int32 // East-West
	EHN      []int32 // North-South
	Checksum [3]uint8
}

type Filter struct {
	a1 float64
	a2 float64
	b0 float64
	b1 float64
	b2 float64
}
