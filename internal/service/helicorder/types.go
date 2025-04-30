package helicorder

import (
	"context"
	"sync"

	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/timesource"
)

const ID = "service_helicorder"

const (
	IMAGE_FORMAT_PNG = "png"
	IMAGE_FORMAT_SVG = "svg"
)

const (
	TIMESPAN_10_MINUTES int64 = 10
	TIMESPAN_15_MINUTES int64 = 15
	TIMESPAN_30_MINUTES int64 = 30
	TIMESPAN_60_MINUTES int64 = 60
)

type HelicorderServiceImpl struct {
	mu     sync.Mutex
	status service.Status

	wg       sync.WaitGroup
	ctx      context.Context
	cancelFn context.CancelFunc

	timeSource   *timesource.Source
	hardwareDev  hardware.IHardware
	dataProvider provider
	channelCodes []string

	filePath    string
	imageFormat string

	timeSpan    int
	lifeCycle   int
	imageSize   int
	spanSamples int
	lineWidth   float64
	scaleFactor float64
}
