package miniseed

import (
	"github.com/anyshake/observer/drivers/explorer"
	"github.com/bclswl0827/mseedio"
)

const (
	MINISEED_BIT_ORDER        = mseedio.MSBFIRST
	MINISEED_ENCODE_TYPE      = mseedio.STEIM2
	MINISEED_WRITE_INTERVAL   = 5
	MINISEED_CLEANUP_INTERVAL = 60
)

type MiniSeedService struct {
	miniseedBuffer       []explorer.ExplorerData
	miniseedSequence     map[string]int // Indepedent sequence number for Z, E, N
	writeBufferInterval  int
	writeBufferCountDown int
	cleanUpCountDown     int
	lifeCycle            int
	basePath             string
	stationCode          string
	networkCode          string
	locationCode         string
	channelPrefix        string
}
