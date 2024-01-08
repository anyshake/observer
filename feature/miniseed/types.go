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
	DataBuffer []int32
	Samples    int32
	SeqNum     int64
}

type miniSEEDBuffer struct {
	TimeStamp     time.Time
	ChannelBuffer map[string]*channelBuffer
}
