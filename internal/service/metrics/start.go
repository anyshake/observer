package metrics

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/system"
	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	_trace "go.opentelemetry.io/otel/trace"
)

func (s *MetricsServiceImpl) getBuildInfo() []attribute.KeyValue {
	buildTime := s.build.GetTime().Format(time.RFC3339)
	buildChannel := s.build.GetChannel()
	buildCommit := s.build.GetCommit()
	toolchainId := s.build.GetToolchainId()
	toolchainObj := s.build.GetToolchain()

	values := []attribute.KeyValue{
		attribute.String("service.build.time", buildTime),
		attribute.String("service.build.channel", buildChannel),
		attribute.String("service.build.commit", buildCommit),
		attribute.String("service.build.toolchainId", toolchainId),
	}
	if toolchainObj != nil {
		extraValues := []attribute.KeyValue{
			attribute.String("service.build.toolchain.name", toolchainObj.Name),
			attribute.String("service.build.toolchain.name", toolchainObj.GOOS),
			attribute.String("service.build.toolchain.name", toolchainObj.GOARCH),
			attribute.String("service.build.toolchain.name", toolchainObj.GOARM),
			attribute.String("service.build.toolchain.name", toolchainObj.GOMIPS),
		}
		values = append(values, extraValues...)
	}
	return values
}

func (s *MetricsServiceImpl) Start(ctx context.Context) error {
	if err := s.lifecycle.Lock(ctx); err != nil {
		return err
	}
	defer s.lifecycle.Unlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.GetIsRunning() {
		return errors.New("metrics service is already running")
	}
	if err := s.parentCtx.Err(); err != nil {
		return err
	}
	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(s.parentCtx)
	}

	otel.SetLogger(logr.Discard())

	oltpCtx, oltpCancel := context.WithCancel(s.parentCtx)
	exporter, err := otlptracehttp.New(
		oltpCtx,
		otlptracehttp.WithEndpoint(OTLP_EXPORTER_ENDPOINT),
		otlptracehttp.WithTimeout(METRICS_REPORT_TIMEOUT),
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{Enabled: false, MaxElapsedTime: 0}),
	)
	if err != nil {
		oltpCancel()
		return fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	hashedHostname := md5.Sum([]byte(s.hostInfo.Hostname))
	resources := []attribute.KeyValue{
		semconv.TelemetrySDKLanguageGo,
		semconv.ServiceNameKey.String(OTLP_SERVICE_NAME),
		semconv.ServiceVersionKey.String(s.version.String()),
		semconv.ProcessRuntimeVersionKey.String(runtime.Version()),
		semconv.OSVersionKey.String(s.hostInfo.PlatformVersion),
		semconv.OSNameKey.String(s.hostInfo.Platform),
		semconv.OSTypeKey.String(runtime.GOOS),
		semconv.HostArchKey.String(runtime.GOARCH),
		// Report only hashed hostname for privacy
		semconv.HostNameKey.String(hex.EncodeToString(hashedHostname[:])),
	}
	resources = append(resources, s.getBuildInfo()...)

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewSchemaless(resources...)),
	)
	otel.SetTracerProvider(tracerProvider)

	runCtx, runCancel := s.ctx, s.cancelFn
	s.status.SetStartedAt(s.timeSource.Now())
	s.status.SetIsRunning(true)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(METRICS_REPORT_INTERVAL)
		tracer := otel.Tracer(OTLP_TRACER_NAME)
		defer func() {
			ticker.Stop()
			runCancel()
			shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(runCtx), METRICS_REPORT_TIMEOUT)
			if err := tracerProvider.Shutdown(shutdownCtx); err != nil {
				logger.GetLogger(ID).Errorf("failed to shutdown OTLP tracer provider: %v", err)
			}
			cancel()
			oltpCancel()
			s.status.SetStoppedAt(s.timeSource.Now())
			s.status.SetIsRunning(false)
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectedly crashed, recovered from panic: %v\n%s", r, debug.Stack())
			}
		}()

		for {
			select {
			case <-runCtx.Done():
				return
			case <-ticker.C:
				if err := s.checkConnectivity(runCtx); err == nil {
					s.reportCurrentStatus(runCtx, tracer)
				}
			}
		}
	}()

	return nil
}

func (s *MetricsServiceImpl) checkConnectivity(ctx context.Context) error {
	client := &http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}}}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, CONNECTIVITY_CHECK_URL, nil)
	if err != nil {
		return fmt.Errorf("pre-check failed: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("pre-check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *MetricsServiceImpl) reportCurrentStatus(ctx context.Context, tracer _trace.Tracer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithTimeout(ctx, METRICS_REPORT_TIMEOUT)
	defer cancel()

	currentTime := s.timeSource.Now()
	systemTime := time.Now().UTC()
	appElapsed := currentTime.Sub(s.startTime)

	ctx, appSpan := tracer.Start(ctx, "app-status")
	appSpan.SetAttributes(attribute.String("app.current_time", currentTime.Format(time.RFC3339)))
	appSpan.SetAttributes(attribute.String("app.system_time", systemTime.Format(time.RFC3339)))
	appSpan.SetAttributes(attribute.String("app.start_time", s.startTime.Format(time.RFC3339)))
	appSpan.SetAttributes(attribute.String("app.elapsed", appElapsed.String()))
	defer appSpan.End()

	var errObj error

	cpuPercent, err := system.GetCpuPercent()
	if err != nil {
		errObj = errors.Join(errObj, err)
	}
	memPercent, err := system.GetMemoryPercent()
	if err != nil {
		errObj = errors.Join(errObj, err)
	}
	diskPercent, err := system.GetDiskPercent()
	if err != nil {
		errObj = errors.Join(errObj, err)
	}
	osUptime, err := system.GetOsUptime()
	if err != nil {
		errObj = errors.Join(errObj, err)
	}

	if errObj != nil {
		appSpan.SetStatus(codes.Error, "error occurred in system metrics")
	} else {
		appSpan.SetStatus(codes.Ok, "successfully got status for entire application")
	}

	_, systemSpan := tracer.Start(ctx, "system", _trace.WithSpanKind(_trace.SpanKindInternal))
	if errObj != nil {
		systemSpan.SetStatus(codes.Error, fmt.Sprintf("error occurred while retrieving system metrics: %v", errObj))
	} else {
		systemSpan.SetStatus(codes.Ok, "successfully got current system status")
		systemSpan.SetAttributes(attribute.String("system.os_uptime", strconv.FormatInt(osUptime, 10)))
		systemSpan.SetAttributes(attribute.Float64("system.cpu_usage", cpuPercent))
		systemSpan.SetAttributes(attribute.Float64("system.memory_usage", memPercent))
		systemSpan.SetAttributes(attribute.Float64("system.disk_usage", diskPercent))
	}
	systemSpan.End()

	deviceConfig := s.hardwareDev.GetConfig()
	deviceStatus := s.hardwareDev.GetStatus()
	_, hardwareSpan := tracer.Start(ctx, "hardware", _trace.WithSpanKind(_trace.SpanKindInternal))
	hardwareSpan.SetStatus(codes.Ok, "successfully got current hardware status")
	hardwareSpan.SetAttributes(attribute.String("hardware.model", deviceConfig.GetModel()))
	hardwareSpan.SetAttributes(attribute.String("hardware.device_id", s.hardwareDev.GetDeviceId()))
	hardwareSpan.SetAttributes(attribute.String("hardware.sample_rate", strconv.Itoa(deviceConfig.GetSampleRate())))
	hardwareSpan.SetAttributes(attribute.String("hardware.errors", strconv.FormatInt(deviceStatus.GetErrors(), 10)))
	hardwareSpan.SetAttributes(attribute.String("hardware.frames", strconv.FormatInt(deviceStatus.GetFrames(), 10)))
	hardwareSpan.SetAttributes(attribute.String("hardware.messages", strconv.FormatInt(deviceStatus.GetMessages(), 10)))
	hardwareSpan.SetAttributes(attribute.String("hardware.protocol", deviceConfig.GetProtocol()))
	hardwareSpan.SetAttributes(attribute.Bool("hardware.gnss_available", deviceConfig.GetGnssAvailability()))
	hardwareSpan.SetAttributes(attribute.StringSlice("hardware.channels", deviceConfig.GetChannelCodes()))
	hardwareSpan.End()
}
