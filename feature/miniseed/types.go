package miniseed

import (
	"github.com/bclswl0827/mseedio"
)

const MODULE string = "miniseed"

const (
	MAX_DURATION  float64 = 3.0
	BIT_ORDER     int     = mseedio.MSBFIRST
	ENCODING_TYPE int     = mseedio.STEIM2
)

type MiniSEED struct{}
