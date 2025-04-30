package miniseed

import (
	"context"
	"sync"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/timesource"
)

const ID = "service_miniseed"

const (
	MINISEED_APPEND_INTERVAL  = 10
	MINISEED_CLEANUP_INTERVAL = 3600
)

type buffer struct {
	SampleRate int
	Timestamp  int64
	explorer.ChannelData
}

type MiniSeedServiceImpl struct {
	mu     sync.Mutex
	status service.Status

	wg       sync.WaitGroup
	ctx      context.Context
	cancelFn context.CancelFunc

	hardwareDev   hardware.IHardware
	timeSource    *timesource.Source
	actionHandler *action.Handler

	stationCode  string
	networkCode  string
	locationCode string
	filePath     string
	lifeCycle    int
	useCompress  bool

	cleanupCountDown int
	appendCountDown  int
	dataSequence     sequence
	recordBuffer     [][]buffer
}
