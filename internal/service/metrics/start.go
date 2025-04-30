package metrics

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"time"

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

func (s *MetricsServiceImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(context.Background())
	}

	s.oltpCtx = context.Background()
	exporter, err := otlptracehttp.New(
		s.oltpCtx,
		otlptracehttp.WithEndpoint(OTLP_EXPORTER_ENDPOINT),
		otlptracehttp.WithTimeout(METRICS_REPORT_TIMEOUT),
	)
	if err != nil {
		return fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	tp := trace.NewTracerProvider(
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
	s.oltpTracerProvider = tp

	otel.SetLogger(logr.Discard())
	otel.SetTracerProvider(tp)
	tracer := otel.Tracer(OTLP_TRACER_NAME)

	go func() {
		s.status.SetStartedAt(s.timeSource.Get())
		s.status.SetIsRunning(true)

		timer := time.NewTimer(30 * time.Second)

		for {
			timer.Reset(METRICS_REPORT_INTERVAL)

			select {
			case <-s.ctx.Done():
				timer.Stop()
				s.wg.Done()
				return
			case <-timer.C:
				s.reportCurrentStatus(tracer)
			}
		}
	}()

	s.wg.Add(1)
	return nil
}

func (s *MetricsServiceImpl) reportCurrentStatus(tracer _trace.Tracer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx := context.Background()
	currentTime := s.timeSource.Get()
	appElapsed := currentTime.Sub(s.startTime)

	ctx, appSpan := tracer.Start(ctx, "app-status")
	appSpan.SetAttributes(attribute.String("app.current_time", currentTime.Format(time.RFC3339)))
	appSpan.SetAttributes(attribute.String("app.start_time", s.startTime.Format(time.RFC3339)))
	appSpan.SetAttributes(attribute.String("app.elapsed", appElapsed.String()))
	defer appSpan.End()

	var errData error

	cpuModel, err := system.GetCpuModel()
	if err != nil {
		errData = errors.Join(errData, err)
	}
	cpuPercent, err := system.GetCpuPercent()
	if err != nil {
		errData = errors.Join(errData, err)
	}
	memPercent, err := system.GetMemoryPercent()
	if err != nil {
		errData = errors.Join(errData, err)
	}
	diskPercent, err := system.GetDiskPercent()
	if err != nil {
		errData = errors.Join(errData, err)
	}
	osUptime, err := system.GetOsUptime()
	if err != nil {
		errData = errors.Join(errData, err)
	}

	if errData != nil {
		appSpan.SetStatus(codes.Error, "error occurred in system metrics")
	} else {
		appSpan.SetStatus(codes.Ok, "successfully got status for entire application")
	}

	_, systemSpan := tracer.Start(ctx, "system", _trace.WithSpanKind(_trace.SpanKindInternal))
	if errData != nil {
		systemSpan.SetStatus(codes.Error, fmt.Sprintf("error occurred while retrieving system metrics: %s", err.Error()))
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
