package miniseed

import (
	"time"

	"github.com/bclswl0827/mseedio"
)

const MODULE string = "miniseed"

const (
	MAX_DURATION  float64 = 4.0
	BIT_ORDER     int     = mseedio.MSBFIRST
	ENCODING_TYPE int     = mseedio.STEIM2
)

type MiniSEED struct{}

type channelBuffer struct {
	SampleRate int32
	DataBuffer []int32
}

type miniSEEDBuffer struct {
	SeqNum    int64
	TimeStamp time.Time
	EHZ       *channelBuffer
	EHE       *channelBuffer
	EHN       *channelBuffer
}
