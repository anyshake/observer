package quakesense

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"runtime/debug"
	"time"

	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/ringbuf"
	"github.com/bclswl0827/eewgo"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/samber/lo"
)

func (s *QuakeSenseServiceImpl) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ctx.Err() != nil {
		s.ctx, s.cancelFn = context.WithCancel(context.Background())
	}

	mqttClientOptions := mqtt.NewClientOptions()
	mqttClientOptions.AddBroker(s.mqttBroker)
	mqttClientOptions.SetAutoReconnect(true)
	mqttClientOptions.SetKeepAlive(30 * time.Second)
	mqttClientOptions.SetClientID("anyshake-observer")
	mqttClientOptions.SetConnectTimeout(10 * time.Second)
	if s.mqttUsername != "" && s.mqttPassword != "" {
		mqttClientOptions.SetUsername(s.mqttUsername)
		mqttClientOptions.SetPassword(s.mqttPassword)
	}
	mqttClientOptions.OnReconnecting = func(c mqtt.Client, options *mqtt.ClientOptions) {
		logger.GetLogger(ID).Warnf("reconnecting to MQTT broker: %s", s.mqttBroker)
	}
	mqttClientOptions.OnConnect = func(c mqtt.Client) {
		logger.GetLogger(ID).Infof("connected to MQTT broker: %s", s.mqttBroker)
	}
	mqttClientOptions.OnConnectionLost = func(c mqtt.Client, err error) {
		logger.GetLogger(ID).Warnf("connection to MQTT broker lost: %v", err)
	}

	s.mqttClient = mqtt.NewClient(mqttClientOptions)
	if token := s.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.GetLogger(ID).Errorf("service unexpectly stopped, recovered from panic: %v\n%s", r, debug.Stack())
				_ = s.Stop()
			}
		}()

		s.status.SetStartedAt(s.timeSource.Get())
		s.status.SetIsRunning(true)

		var lastTriggeredTime time.Time

		s.hardwareDev.Subscribe(ID, func(t time.Time, di *explorer.DeviceConfig, dv *explorer.DeviceVariable, cd []explorer.ChannelData) {
			s.mu.Lock()
			defer s.mu.Unlock()

			targetChannel, targetChannelFound := lo.Find(cd, func(c explorer.ChannelData) bool { return c.ChannelCode == s.monitorChannel })
			if !targetChannelFound {
				logger.GetLogger(ID).Warnf("target monitoring channel %s not found", s.monitorChannel)
				return
			}

			currentSamplerate := di.GetSampleRate()
			bufferSize := int(s.ltaWindow * float64(currentSamplerate))
			if s.prevSamplerate == 0 || s.prevSamplerate != currentSamplerate {
				s.prevSamplerate = currentSamplerate
				s.channelBuffer = ringbuf.New[float64](bufferSize)

				switch s.filterType {
				case "bandpass":
					s.filterKernel = eewgo.NewFilter(eewgo.UseBandPassFilter(s.minFreq, s.maxFreq, float64(s.prevSamplerate), FILTER_NUM_TAPS))
				case "lowpass":
					s.filterKernel = eewgo.NewFilter(eewgo.UseLowPassFilter(s.maxFreq, float64(s.prevSamplerate), FILTER_NUM_TAPS))
				case "highpass":
					s.filterKernel = eewgo.NewFilter(eewgo.UseHighPassFilter(s.minFreq, float64(s.prevSamplerate), FILTER_NUM_TAPS))
				default:
					s.filterKernel = eewgo.NewFilter(eewgo.UseBandPassFilter(0.5, 10, float64(s.prevSamplerate), FILTER_NUM_TAPS))
				}
			}

			channelData := lo.Map(targetChannel.Data, func(v int32, _ int) float64 { return float64(v) })
			s.channelBuffer.Push(channelData...)

			if s.channelBuffer.Len() < bufferSize {
				logger.GetLogger(ID).Infof("waiting for %d samples to fill LTA window", bufferSize)
				return
			}

			filtered := s.filterKernel.Apply(s.channelBuffer.Values())
			var staLtaArr []float64
			switch s.triggerMethod {
			case CLASSIC_STA_LTA:
				staLtaArr = eewgo.ClassicStaLta(filtered, int(s.staWindow*float64(s.prevSamplerate)), int(s.ltaWindow*float64(s.prevSamplerate)))
			// case DELAYED_STA_LTA:
			// 	staLtaArr = eewgo.DelayedStaLta(filtered, int(s.staWindow*float64(s.prevSamplerate)), int(s.ltaWindow*float64(s.prevSamplerate)))
			// case RECURSIVE_STA_LTA:
			// 	staLtaArr = eewgo.RecursiveStaLta(filtered, int(s.staWindow*float64(s.prevSamplerate)), int(s.ltaWindow*float64(s.prevSamplerate)))
			case Z_DETECT:
				staLtaArr = eewgo.ZDetect(filtered, int(s.staWindow*float64(s.prevSamplerate)))
			default:
				logger.GetLogger(ID).Warnf("unknown trigger method sepcified: %s", s.triggerMethod)
				return
			}

			onsets := eewgo.TriggerOnset(staLtaArr, s.trigOn, s.trigOff, math.MaxInt32, false)
			if len(onsets) > 0 {
				if s.throttleSeconds > 0 && !lastTriggeredTime.IsZero() {
					elapsed := t.Sub(lastTriggeredTime)
					if elapsed < time.Duration(s.throttleSeconds)*time.Second {
						return
					}
				}

				logger.GetLogger(ID).Infof("detected %d seismic event at UTC time: %s", len(onsets), t.UTC().Format(time.RFC3339))
				lastTriggeredTime = t

				latitude, longitude, elevation, err := s.hardwareDev.GetCoordinates(true)
				if err != nil {
					logger.GetLogger(ID).Warnf("failed to get coordinates: %v", err)
					return
				}

				startTime := t.Add(-time.Duration(bufferSize/currentSamplerate) * time.Second)
				for _, onset := range onsets {
					triggerTime := startTime.Add(time.Duration(onset[0]) * time.Second / time.Duration(s.prevSamplerate))
					payload, err := json.Marshal(map[string]any{
						"trigger_method":      s.triggerMethod,
						"trigger_time":        triggerTime.UnixMilli(),
						"station_name":        s.stationName,
						"station_description": s.stationDescription,
						"station_country":     s.stationCountry,
						"station_place":       s.stationPlace,
						"station_affiliation": s.stationAffiliation,
						"latitude":            latitude,
						"longitude":           longitude,
						"elevation":           elevation,
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
						"num_taps":            FILTER_NUM_TAPS,
						"sample_rate":         s.prevSamplerate,
						"channel_code":        s.monitorChannel,
					})
					if err != nil {
						logger.GetLogger(ID).Errorf("failed to marshal payload: %v", err)
						return
					}
					token := s.mqttClient.Publish(s.mqttTopic, 0, false, string(payload))
					if token.Wait() && token.Error() != nil {
						logger.GetLogger(ID).Errorf("failed to publish MQTT message: %v", token.Error())
					}
				}
			}
		})

		<-s.ctx.Done()
		s.wg.Done()
	}()

	s.wg.Add(1)
	return nil
}
