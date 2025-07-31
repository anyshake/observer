package metrics

import (
	"context"
	"crypto/tls"
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

func (s *MetricsServiceImpl) handleInterrupt(ticker *time.Ticker) {
	ticker.Stop()
	s.wg.Done()
}

func (s *MetricsServiceImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(context.Background())
	}

	otel.SetLogger(logr.Discard())

	s.oltpCtx, s.oltpCtxCancelFn = context.WithCancel(context.Background())
	exporter, err := otlptracehttp.New(
		s.oltpCtx,
		otlptracehttp.WithEndpoint(OTLP_EXPORTER_ENDPOINT),
		otlptracehttp.WithTimeout(METRICS_REPORT_TIMEOUT),
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{Enabled: false, MaxElapsedTime: 0}),
	)
	if err != nil {
		return fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	s.oltpTracerProvider = trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewSchemaless(
			semconv.TelemetrySDKLanguageGo,
			semconv.ServiceNameKey.String(OTLP_SERVICE_NAME),
			semconv.ServiceVersionKey.String(s.binaryVersion),
			attribute.String("service.build_platform", s.buildPlatform),
			attribute.String("service.commit_hash", s.commitHash),
			semconv.ProcessRuntimeVersionKey.String(runtime.Version()),
			semconv.HostID(s.hostInfo.HostID),
			semconv.OSVersionKey.String(s.hostInfo.PlatformVersion),
			semconv.OSNameKey.String(s.hostInfo.Platform),
			semconv.HostNameKey.String(s.hostInfo.Hostname),
			semconv.OSTypeKey.String(runtime.GOOS),
			semconv.HostArchKey.String(runtime.GOARCH),
		)),
	)
	otel.SetTracerProvider(s.oltpTracerProvider)

	go func() {
		ticker := time.NewTicker(METRICS_REPORT_INTERVAL)
		tracer := otel.Tracer(OTLP_TRACER_NAME)

		s.status.SetStartedAt(s.timeSource.Now())
		s.status.SetIsRunning(true)
		defer func() {
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectly crashed, recovered from panic: %v\n%s", r, debug.Stack())
				s.handleInterrupt(ticker)
				_ = s.Stop()
			}
		}()

		for {
			select {
			case <-s.ctx.Done():
				s.handleInterrupt(ticker)
				return
			case <-ticker.C:
				if err := s.checkConnectivity(); err == nil {
					s.reportCurrentStatus(tracer)
				}
			}
		}
	}()

	s.wg.Add(1)
	return nil
}

func (s *MetricsServiceImpl) checkConnectivity() error {
	client := &http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}}}
	resp, err := client.Get(CONNECTIVITY_CHECK_URL)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *MetricsServiceImpl) reportCurrentStatus(tracer _trace.Tracer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), METRICS_REPORT_TIMEOUT)
	defer cancel()

	currentTime := s.timeSource.Now()
	appElapsed := currentTime.Sub(s.startTime)

	ctx, appSpan := tracer.Start(ctx, "app-status")
	appSpan.SetAttributes(attribute.String("app.current_time", currentTime.Format(time.RFC3339)))
	appSpan.SetAttributes(attribute.String("app.start_time", s.startTime.Format(time.RFC3339)))
	appSpan.SetAttributes(attribute.String("app.elapsed", appElapsed.String()))
	defer appSpan.End()

	var errObj error

	cpuModel, err := system.GetCpuModel()
	if err != nil {
		errObj = errors.Join(errObj, err)
	}
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
		systemSpan.SetAttributes(attribute.String("system.cpu_model", cpuModel))
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
