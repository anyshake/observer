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
)

const (
	CACHE_STORAGE_DISABLED = "disabled"
	CACHE_STORAGE_DISK     = "disk"
	CACHE_STORAGE_MEMORY   = "memory"
	CACHE_DEFAULT_PATH     = "./service_data/helicorder/.helicorder-cache"
)

type HelicorderServiceImpl struct {
	lifecycle service.Lifecycle
	mu        sync.Mutex
	status    service.Status

	wg        sync.WaitGroup
	parentCtx context.Context
	ctx       context.Context
	cancelFn  context.CancelFunc

	timeSource   *timesource.Source
	hardwareDev  hardware.IHardware
	dataProvider provider

	filePath     string
	imageFormat  string
	cacheStorage string
	cachePath    string

	timeSpan     int
	lifeCycle    int
	imageSize    int
	spanSamples  int
	lineWidth    float64
	scaleFactors []float64
}
