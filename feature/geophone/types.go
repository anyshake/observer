package geophone

import "time"

const MODULE string = "geophone"

const (
	// READY_THRESHOLD should be strictly 1 second
	READY_THRESHOLD time.Duration = 1 * time.Second
	// TIMEOUT_THRESHOLD should be greater than READY_THRESHOLD
	TIMEOUT_THRESHOLD time.Duration = 3 * time.Second
)

var (
	// RESET_WORD resets geophone ADC module
	RESET_WORD = [...]byte{0x61}
	// SYNC_WORD indicates a data packet is following
	SYNC_WORD = [...]byte{0xFC, 0x1B}
	// ACK_WORD indicates a valid command is received
	ACK_WORD = [...]byte{0xFC, 0x2B}
)

type Geophone struct {
	Ticker *time.Ticker
}

type Packet struct {
	EHZ      []int32 // Vertical
	EHE      []int32 // East-West
	EHN      []int32 // North-South
	Checksum [3]uint8
}
