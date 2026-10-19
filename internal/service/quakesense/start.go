package quakesense

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"runtime/debug"
	"time"

	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/notification"
	"github.com/anyshake/observer/internal/service"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/ringbuf"
	"github.com/bclswl0827/obsgo/signal"
	"github.com/bclswl0827/obsgo/trigger"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	quakeSenseAnalysisInterval = 250 * time.Millisecond
	quakeSenseConnectTimeout   = 10 * time.Second
	quakeSensePublishTimeout   = 2 * time.Second
	quakeSenseWarningInterval  = 10 * time.Second
)

func quakeSenseNotificationMessage(eventTime time.Time, eventCount int) string {
	return fmt.Sprintf(
		"[%s UTC] QuakeSense detected %d seismic event(s).",
		eventTime.UTC().Format("2006-01-02 15:04:05"),
		eventCount,
	)
}

func waitMQTTToken(ctx context.Context, token mqtt.Token, timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-token.Done():
		return true
	case <-ctx.Done():
		return false
	case <-timer.C:
		return false
	}
}

func (s *QuakeSenseServiceImpl) filterBufferSize(sampleRate int) int {
	ltaSamples := max(1, int(s.ltaWindow*float64(sampleRate)))
	if s.triggerMethod == DELAYED_STA_LTA {
		// obsgo's delayed STA/LTA uses samples from STA+LTA+50 samples
		// before the current position.
		staSamples := max(1, int(s.staWindow*float64(sampleRate)))
		return max(ltaSamples, staSamples+ltaSamples+51)
	}
	return ltaSamples
}

func (s *QuakeSenseServiceImpl) newFilter(sampleRate int) (*signal.IIRFilter, error) {
	if sampleRate <= 0 {
		return nil, fmt.Errorf("sample rate must be positive")
	}

	samplePeriod := 1 / float64(sampleRate)
	switch s.filterType {
	case NO_FILTER:
		return nil, nil
	case BAND_PASS_FILTER:
		return signal.NewButterworth(FILTER_ORDER, signal.BandPass, s.minFreq, s.maxFreq, samplePeriod)
	case LOW_PASS_FILTER:
		return signal.NewButterworth(FILTER_ORDER, signal.LowPass, 0, s.maxFreq, samplePeriod)
	case HIGH_PASS_FILTER:
		return signal.NewButterworth(FILTER_ORDER, signal.HighPass, s.minFreq, 0, samplePeriod)
	default:
		return nil, fmt.Errorf("unknown filter type %q", s.filterType)
	}
}

func (s *QuakeSenseServiceImpl) Start(ctx context.Context) error {
	if err := s.lifecycle.Lock(ctx); err != nil {
		return err
	}
	defer s.lifecycle.Unlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.GetIsRunning() {
		return errors.New("QuakeSense service is already running")
	}
	if err := s.parentCtx.Err(); err != nil {
		return err
	}
	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(s.parentCtx)
	}

	mqttBroker := s.mqttBroker
	mqttClientOptions := mqtt.NewClientOptions()
	mqttClientOptions.AddBroker(mqttBroker)
	mqttClientOptions.SetAutoReconnect(true)
	mqttClientOptions.SetKeepAlive(30 * time.Second)
	mqttClientOptions.SetClientID(s.mqttClientId)
	mqttClientOptions.SetConnectTimeout(quakeSenseConnectTimeout)
	mqttClientOptions.SetWriteTimeout(quakeSensePublishTimeout)
	if s.mqttUsername != "" {
		mqttClientOptions.SetUsername(s.mqttUsername)
	}
	if s.mqttPassword != "" {
		mqttClientOptions.SetPassword(s.mqttPassword)
	}
	mqttClientOptions.OnReconnecting = func(c mqtt.Client, options *mqtt.ClientOptions) {
		logger.GetLogger(ID).Warnf("reconnecting to MQTT broker: %s", mqttBroker)
	}
	mqttClientOptions.OnConnect = func(c mqtt.Client) {
		logger.GetLogger(ID).Infof("connected to MQTT broker: %s", mqttBroker)
	}
	mqttClientOptions.OnConnectionLost = func(c mqtt.Client, err error) {
		logger.GetLogger(ID).Warnf("connection to MQTT broker lost: %v", err)
	}

	mqttClient := mqtt.NewClient(mqttClientOptions)
	token := mqttClient.Connect()
	if !waitMQTTToken(ctx, token, quakeSenseConnectTimeout) {
		mqttClient.Disconnect(0)
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("failed to connect to MQTT broker: %w", err)
		}
		return fmt.Errorf("failed to connect to MQTT broker: timed out after %s", quakeSenseConnectTimeout)
	}
	if err := token.Error(); err != nil {
		mqttClient.Disconnect(0)
		return fmt.Errorf("failed to connect to MQTT broker: %w", err)
	}

	runCtx, runCancel := s.ctx, s.cancelFn
	s.status.SetStartedAt(s.timeSource.Now())
	s.status.SetIsRunning(true)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer func() {
			runCancel()
			_ = s.hardwareDev.UnsubscribeRealtime(ID)
			mqttClient.Disconnect(100)
			s.mu.Lock()
			s.prevSamplerate = 0
			s.filterKernel = nil
			if s.channelBuffer != nil {
				s.channelBuffer.Reset()
			}
			s.mu.Unlock()
			s.status.SetStoppedAt(s.timeSource.Now())
			s.status.SetIsRunning(false)
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectedly crashed, recovered from panic: %v\n%s", r, debug.Stack())
			}
		}()

		var lastTriggeredAt time.Time
		var lastTriggerTime time.Time
		var lastStreamWarning time.Time
		var pendingStart time.Time
		var expectedNextTimestamp time.Time
		pendingSamples := make([]int32, 0)
		logStreamWarning := func(format string, args ...any) {
			now := s.timeSource.Now()
			if !lastStreamWarning.IsZero() && now.Sub(lastStreamWarning) < quakeSenseWarningInterval {
				return
			}
			lastStreamWarning = now
			logger.GetLogger(ID).Warnf(format, args...)
		}
		type detection struct {
			triggerTimes []time.Time
			topic        string
			payload      map[string]any
			notify       bool
		}
		detect := func(sensorEvent explorer.Event) *detection {
			s.mu.Lock()
			defer s.mu.Unlock()
			if runCtx.Err() != nil {
				return nil
			}

			var targetChannel *explorer.ChannelData
			for i := range sensorEvent.ChannelData {
				if sensorEvent.ChannelData[i].ChannelCode == s.monitorChannel {
					targetChannel = &sensorEvent.ChannelData[i]
					break
				}
			}
			if targetChannel == nil {
				logStreamWarning("target monitoring channel %s not found", s.monitorChannel)
				return nil
			}
			if len(targetChannel.Data) == 0 {
				return nil
			}

			currentSampleRate := sensorEvent.SampleRate
			if currentSampleRate <= 0 {
				logStreamWarning("invalid device sample rate: %d", currentSampleRate)
				return nil
			}

			bufferSize := s.filterBufferSize(currentSampleRate)
			if s.prevSamplerate != currentSampleRate || s.channelBuffer == nil {
				filterKernel, err := s.newFilter(currentSampleRate)
				if err != nil {
					s.prevSamplerate = 0
					logger.GetLogger(ID).Warnf("failed to create %s filter for sample rate %d: %v", s.filterType, currentSampleRate, err)
					return nil
				}
				s.prevSamplerate = currentSampleRate
				s.filterKernel = filterKernel
				s.channelBuffer = ringbuf.New[float64](bufferSize)
				pendingSamples = pendingSamples[:0]
				pendingStart = time.Time{}
				expectedNextTimestamp = time.Time{}
			}

			samplePeriod := time.Second / time.Duration(currentSampleRate)
			allowedJitter := time.Duration(explorer.ALLOWED_JITTER_MS_NTP) * time.Millisecond
			if sensorEvent.GNSSAvailable {
				allowedJitter = time.Duration(explorer.ALLOWED_JITTER_MS_GNSS) * time.Millisecond
			}
			allowedJitter = max(allowedJitter, 2*samplePeriod)
			if !expectedNextTimestamp.IsZero() {
				timestampDelta := sensorEvent.Timestamp.Sub(expectedNextTimestamp)
				if timestampDelta < -allowedJitter || timestampDelta > allowedJitter {
					s.channelBuffer.Reset()
					if s.filterKernel != nil {
						s.filterKernel.Reset()
					}
					pendingSamples = pendingSamples[:0]
					pendingStart = time.Time{}
					logStreamWarning("realtime stream discontinuity detected, signal state reset: %s", timestampDelta)
				}
			}

			if len(pendingSamples) == 0 {
				pendingStart = sensorEvent.Timestamp
			}
			pendingSamples = append(pendingSamples, targetChannel.Data...)
			expectedNextTimestamp = sensorEvent.Timestamp.Add(
				time.Duration(len(targetChannel.Data)) * time.Second / time.Duration(currentSampleRate),
			)

			analysisSamples := max(1, int(math.Ceil(quakeSenseAnalysisInterval.Seconds()*float64(currentSampleRate))))
			if len(pendingSamples) < analysisSamples {
				return nil
			}

			batchStart := pendingStart
			batchSampleCount := len(pendingSamples)
			channelData := make([]float64, batchSampleCount)
			for i, sample := range pendingSamples {
				channelData[i] = float64(sample)
			}
			pendingSamples = pendingSamples[:0]
			pendingStart = time.Time{}

			filtered := channelData
			if s.filterKernel != nil {
				var err error
				filtered, err = signal.ApplyFilter(s.filterKernel, channelData, 1/float64(currentSampleRate))
				if err != nil {
					logger.GetLogger(ID).Warnf("failed to apply %s filter: %v", s.filterType, err)
					return nil
				}
			}
			s.channelBuffer.Push(filtered...)
			if s.channelBuffer.Len() < bufferSize {
				return nil
			}

			values := s.channelBuffer.Values()
			staSamples := max(1, int(s.staWindow*float64(currentSampleRate)))
			ltaSamples := max(staSamples+1, int(s.ltaWindow*float64(currentSampleRate)))
			var staLta []float64
			switch s.triggerMethod {
			case CLASSIC_STA_LTA:
				staLta = trigger.ClassicStaLta(values, staSamples, ltaSamples)
			case DELAYED_STA_LTA:
				staLta = trigger.DelayedStaLta(values, staSamples, ltaSamples)
			case RECURSIVE_STA_LTA:
				staLta = trigger.RecursiveStaLta(values, staSamples, ltaSamples)
			case Z_DETECT:
				staLta = trigger.ZDetect(values, staSamples)
			default:
				logger.GetLogger(ID).Warnf("unknown trigger method specified: %s", s.triggerMethod)
				return nil
			}

			onsets := trigger.TriggerOnset(staLta, s.trigOn, s.trigOff, math.MaxInt32, false)
			if len(onsets) == 0 {
				return nil
			}

			bufferStart := batchStart.Add(time.Duration(batchSampleCount-len(values)) * time.Second / time.Duration(currentSampleRate))
			triggerTimes := make([]time.Time, 0, len(onsets))
			for _, onset := range onsets {
				if len(onset) == 0 || onset[0] < 0 || onset[0] >= len(values) {
					continue
				}
				triggerTime := bufferStart.Add(time.Duration(onset[0]) * time.Second / time.Duration(currentSampleRate))
				if !lastTriggerTime.IsZero() && !triggerTime.After(lastTriggerTime) {
					continue
				}
				triggerTimes = append(triggerTimes, triggerTime)
			}
			if len(triggerTimes) == 0 {
				return nil
			}
			if s.throttleSeconds > 0 && !lastTriggeredAt.IsZero() && triggerTimes[0].Sub(lastTriggeredAt) < time.Duration(s.throttleSeconds)*time.Second {
				return nil
			}

			return &detection{
				triggerTimes: triggerTimes,
				topic:        s.mqttTopic,
				payload: map[string]any{
					"trigger_method":      s.triggerMethod,
					"station_name":        s.stationName,
					"station_description": s.stationDescription,
					"station_country":     s.stationCountry,
					"station_place":       s.stationPlace,
					"station_affiliation": s.stationAffiliation,
					"station_code":        s.stationCode,
					"network_code":        s.networkCode,
					"location_code":       s.locationCode,
					"sta_window":          s.staWindow,
					"lta_window":          s.ltaWindow,
					"trig_on":             s.trigOn,
					"trig_off":            s.trigOff,
					"filter_type":         s.filterType,
					"min_freq":            s.minFreq,
					"max_freq":            s.maxFreq,
					"filter_order":        FILTER_ORDER,
					"sample_rate":         currentSampleRate,
					"channel_code":        s.monitorChannel,
				},
				notify: s.webNotificationEnabled && s.notifications != nil,
			}
		}
		handler := func(sensorEvent explorer.Event) {
			defer service.RecoverCallback(ID, runCancel)
			detected := detect(sensorEvent)
			if detected == nil {
				return
			}

			latitude, longitude, elevation, err := s.hardwareDev.GetCoordinates(true)
			if err != nil {
				logger.GetLogger(ID).Warnf("failed to get coordinates: %v", err)
				return
			}
			detected.payload["latitude"] = latitude
			detected.payload["longitude"] = longitude
			detected.payload["elevation"] = elevation

			publishedEvents := 0
			var lastPublishedAt time.Time
			for _, triggerTime := range detected.triggerTimes {
				detected.payload["trigger_time"] = triggerTime.UnixMilli()
				payload, err := json.Marshal(detected.payload)
				if err != nil {
					logger.GetLogger(ID).Errorf("failed to marshal payload: %v", err)
					continue
				}
				token := mqttClient.Publish(detected.topic, 0, false, string(payload))
				if !waitMQTTToken(runCtx, token, quakeSensePublishTimeout) {
					if runCtx.Err() != nil {
						return
					}
					logger.GetLogger(ID).Errorf("timed out publishing MQTT message after %s", quakeSensePublishTimeout)
					continue
				}
				if token.Error() != nil {
					logger.GetLogger(ID).Errorf("failed to publish MQTT message: %v", token.Error())
					continue
				}
				lastTriggerTime = triggerTime
				lastPublishedAt = triggerTime
				publishedEvents++
			}
			if publishedEvents == 0 {
				return
			}
			lastTriggeredAt = lastPublishedAt
			logger.GetLogger(ID).Infof("detected %d seismic event at UTC time: %s", publishedEvents, lastPublishedAt.UTC().Format(time.RFC3339))
			if detected.notify {
				s.notifications.Publish(notification.Event{
					ServiceID:  ID,
					Message:    quakeSenseNotificationMessage(lastPublishedAt, publishedEvents),
					Level:      notification.LevelWarning,
					OccurredAt: lastPublishedAt.UnixMilli(),
				})
			}
		}

		if err := s.hardwareDev.SubscribeRealtime(ID, handler); err != nil {
			logger.GetLogger(ID).Errorf("failed to subscribe to realtime hardware message bus: %v", err)
			return
		}

		<-runCtx.Done()
	}()

	return nil
}
