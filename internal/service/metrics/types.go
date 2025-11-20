package metrics

import (
	"context"
	"sync"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/semver"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/anyshake/observer/pkg/unibuild"
	"github.com/shirou/gopsutil/v4/host"
	"go.opentelemetry.io/otel/sdk/trace"
)

const ID = "service_metrics"

const (
	CONNECTIVITY_CHECK_URL  = "https://metrics.anyshake.org"
	OTLP_EXPORTER_ENDPOINT  = "metrics.anyshake.org:443"
	OTLP_SERVICE_NAME       = "anyshake-observer"
	OTLP_TRACER_NAME        = "anyshake-observer-tracer"
	ATTRIBUTE_DUMMY_VALUE   = "<dummy>"
	METRICS_REPORT_INTERVAL = 10 * time.Minute
	METRICS_REPORT_TIMEOUT  = 5 * time.Second
)

type MetricsServiceImpl struct {
	mu     sync.Mutex
	status service.Status

	wg       sync.WaitGroup
	ctx      context.Context
	cancelFn context.CancelFunc

	hardwareDev   hardware.IHardware
	timeSource    *timesource.Source
	actionHandler *action.Handler

	oltpCtx            context.Context
	oltpCtxCancelFn    context.CancelFunc
	oltpTracerProvider *trace.TracerProvider

	startTime time.Time
	hostInfo  *host.InfoStat

	version *semver.Version
	build   *unibuild.UniBuild
}
