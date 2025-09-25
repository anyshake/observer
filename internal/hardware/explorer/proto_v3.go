package explorer

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/anyshake/observer/pkg/fifo"
	"github.com/anyshake/observer/pkg/message"
	"github.com/anyshake/observer/pkg/metadata"
	"github.com/anyshake/observer/pkg/ntpclient"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/anyshake/observer/pkg/transport"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type explorerProtocolPacketV3 struct {
	dataTime  time.Time
	dataBytes []byte
}

type ExplorerProtoImplV3 struct {
	ChannelCodes    []string
	ExplorerOptions ExplorerOptions
	NtpOptions      NtpOptions
	Logger          *logrus.Entry
	TimeSource      *timesource.Source

	Transport  transport.ITransport
	fifoBuffer fifo.Buffer[*explorerProtocolPacketV3]

	// 1 message per second, for archiving service, etc.
	messageBus message.Bus[EventHandler]
	// 1 message per packet, for realtime purposes
	messageBusRealtime message.Bus[EventHandler]

	prevMcuTimestamp    int64
	isDataStreamStable  bool
	timeCalibrationChan chan [2]time.Time

	flagMutex        sync.Mutex
	variableAllSet   bool
	collectedSamples int
	packetTimeObj    time.Time

	deviceStatus   DeviceStatus
	deviceConfig   DeviceConfig
	deviceVariable DeviceVariable
	channelDataBuf []ChannelData
}

func (g *ExplorerProtoImplV3) resetFlags() {
	g.flagMutex.Lock()
	defer g.flagMutex.Unlock()

	g.channelDataBuf = []ChannelData{}
	g.packetTimeObj = time.Time{}
	g.collectedSamples = 0
}

func (g *ExplorerProtoImplV3) resetVariables() {
	g.flagMutex.Lock()
	g.variableAllSet = false
	g.flagMutex.Unlock()

	g.deviceVariable.Reset()
}

func (g *ExplorerProtoImplV3) parsePacketInterval(deviceConfig uint32) time.Duration {
	DATA_PACKET_PACKET_INTERVAL := []int{100, 200, 500, 1000}
	return time.Duration(DATA_PACKET_PACKET_INTERVAL[(deviceConfig>>30)&0x3]) * time.Millisecond
}

func (g *ExplorerProtoImplV3) parseSampleRate(deviceConfig uint32) int {
	DATA_PACKET_SAMPLE_RATES := []int{10, 20, 50, 100, 200, 250, 500, 1000}
	return DATA_PACKET_SAMPLE_RATES[(deviceConfig>>27)&0x7]
}

func (g *ExplorerProtoImplV3) parseGnssAvailibility(deviceConfig uint32) bool {
	return ((deviceConfig >> 26) & 0x1) == 1
}

func (g *ExplorerProtoImplV3) getChannelSize(deviceConfig uint32) (channelChunkLength, totalChannelSize int, channelData []*ChannelData) {
	DATA_PACKET_CHANNEL_TYPE := []string{"disabled", "int16", "int24", "int32"}
	packetInterval := g.parsePacketInterval(deviceConfig)
	sampleRate := g.parseSampleRate(deviceConfig)
	channelChunkLength = int(packetInterval.Milliseconds()) / (1000 / sampleRate)

	for i := 0; i < 8; i++ {
		configVal := (deviceConfig >> (24 - i*2)) & 0x3
		if configVal != 0 {
			byteSize := int(configVal) + 1
			ch := &ChannelData{
				ChannelId: i + 1, // Channel ID starts from 1
				ByteSize:  byteSize,
				DataType:  DATA_PACKET_CHANNEL_TYPE[configVal],
			}
			ch.ChannelCode = fmt.Sprintf("CH%d", ch.ChannelId)
			if i < len(g.ChannelCodes) {
				ch.ChannelCode = g.ChannelCodes[i]
			}
			channelData = append(channelData, ch)

			totalChannelSize += (channelChunkLength * byteSize)
		}
	}

	return channelChunkLength, totalChannelSize, channelData
}

func (g *ExplorerProtoImplV3) getVariableData(mcuTimestamp int64, deviceConfig uint32, variableBytes [4]byte) {
	variableData := binary.LittleEndian.Uint32(variableBytes[:])
	variableBits := deviceConfig & 0x3FF

	switch (mcuTimestamp / 1000) % 10 {
	case 0:
		if variableBits&0x1 != 0 {
			g.deviceVariable.SetDeviceId(&variableData)
		} else {
			g.deviceVariable.SetDeviceId(nil)
		}
	case 1:
		if (variableBits>>1)&0x1 != 0 {
			n := float64(math.Float32frombits(variableData))
			g.deviceVariable.SetLatitude(&n)
		} else {
			g.deviceVariable.SetLatitude(&g.ExplorerOptions.Latitude)
		}
	case 2:
		if (variableBits>>2)&0x1 != 0 {
			n := float64(math.Float32frombits(variableData))
			g.deviceVariable.SetLongitude(&n)
		} else {
			g.deviceVariable.SetLongitude(&g.ExplorerOptions.Longitude)
		}
	case 3:
		if (variableBits>>3)&0x1 != 0 {
			n := float64(math.Float32frombits(variableData))
			g.deviceVariable.SetElevation(&n)
		} else {
			g.deviceVariable.SetElevation(&g.ExplorerOptions.Elevation)
		}
	case 4:
		if (variableBits>>4)&0x1 != 0 {
			n := float64(math.Float32frombits(variableData))
			g.deviceVariable.SetTemperature(&n)
		} else {
			g.deviceVariable.SetTemperature(nil)
		}
	}

	variableAllSet := true
	if variableBits&0x1 != 0 {
		if _, err := g.deviceVariable.GetDeviceId(); err != nil {
			variableAllSet = false
		}
	}

	if (variableBits>>1)&0x1 != 0 {
		if _, err := g.deviceVariable.GetLatitude(false); err != nil {
			variableAllSet = false
		}
	}

	if (variableBits>>2)&0x1 != 0 {
		if _, err := g.deviceVariable.GetLongitude(false); err != nil {
			variableAllSet = false
		}
	}

	if (variableBits>>3)&0x1 != 0 {
		if _, err := g.deviceVariable.GetElevation(); err != nil {
			variableAllSet = false
		}
	}

	if (variableBits>>4)&0x1 != 0 {
		if _, err := g.deviceVariable.GetTemperature(); err != nil {
			variableAllSet = false
		}
	}

	g.variableAllSet = variableAllSet
}

func (g *ExplorerProtoImplV3) getChannelData(channelData []*ChannelData, channelDataBytes []byte, channelChunkLength int) {
	offset := 0

	for _, ch := range channelData {
		ch.Data = make([]int32, channelChunkLength)
		for i := 0; i < channelChunkLength; i++ {
			switch ch.DataType {
			case "int16":
				ch.Data[i] = int32(int16(binary.LittleEndian.Uint16(channelDataBytes[offset : offset+2])))
				offset += 2
			case "int24":
				ch.Data[i] = int32(channelDataBytes[offset]) | int32(channelDataBytes[offset+1])<<8 | int32(channelDataBytes[offset+2])<<16
				offset += 3
			case "int32":
				ch.Data[i] = int32(binary.LittleEndian.Uint32(channelDataBytes[offset : offset+4]))
				offset += 4
			}
		}
	}
}

func (g *ExplorerProtoImplV3) verifyChecksum(packetData []byte) error {
	recvChecksum := packetData[len(packetData)-1]
	calcChecksum := uint8(0)
	for _, b := range packetData[:len(packetData)-1] {
		calcChecksum ^= b
	}
	if recvChecksum != calcChecksum {
		return fmt.Errorf("invalid checksum: expected %v, got %v", recvChecksum, calcChecksum)
	}
	return nil
}

func (g *ExplorerProtoImplV3) Open(ctx context.Context) (context.Context, context.CancelFunc, error) {
	if g.Transport == nil {
		return nil, nil, errors.New("transport is not opened")
	}
	if err := g.Transport.Open(); err != nil {
		return nil, nil, fmt.Errorf("failed to open transport: %w", err)
	}
	if err := g.Flush(); err != nil {
		return nil, nil, fmt.Errorf("failed to flush transport: %w", err)
	}

	if g.Logger == nil {
		return nil, nil, errors.New("logger is not set")
	}
	ntpClient, err := ntpclient.New(g.NtpOptions.Pool, g.NtpOptions.Retry, g.NtpOptions.ReadTimeout, timesource.MonotonicNow)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ntp client: %w", err)
	}

	subCtx, cancelFn := context.WithCancel(ctx)

	g.fifoBuffer = fifo.New[*explorerProtocolPacketV3](512)
	g.messageBus = message.NewBus[EventHandler](EXPLORER_STREAM_TOPIC, 1024)
	g.messageBusRealtime = message.NewBus[EventHandler](EXPLORER_REALTIME_STREAM_TOPIC, 1024)
	g.deviceStatus.SetUpdatedAt(time.Unix(0, 0))
	g.deviceConfig.SetProtocol(g.ExplorerOptions.Protocol)
	g.deviceConfig.SetModel(filepath.Base(g.ExplorerOptions.Model))
	g.timeCalibrationChan = make(chan [2]time.Time)

	var initFlag int32
	atomic.StoreInt32(&initFlag, 0)
	readyChan := make(chan struct{})

	DATA_PACKET_HEADER := []byte{0x01, 0xFE}
	DATA_PACKET_TAILER := []byte{0xEF, 0x10}

	go func() {
		timeDiffSamples := make([]int64, 0, STABLE_CHECK_SAMPLES)
		packetDelim := append(DATA_PACKET_TAILER, DATA_PACKET_HEADER...)
		g.isDataStreamStable = false

		for timeSourceInitialized := false; ; {
			select {
			case <-subCtx.Done():
				g.Logger.Info("exiting from data packet reader")
				if atomic.LoadInt32(&initFlag) == 0 {
					close(readyChan)
				}
				return
			default:
			}

			// Assume that the longest packet interval is 1000 ms
			// With 8 channels and 1000 samples per second per channel in int32
			// That would be 8 channels * 4 bytes * (1000 ms / (1000 / 1000 SPS)) for maximum channel data size
			recvBuf, timeout, _, err := g.Transport.ReadUntil(
				subCtx,
				2 /*header*/ +8 /*timestamp*/ +4 /*config*/ +4 /*variable*/ +32000 /*max channel data*/ +1 /*checksum*/ +2, /*tailer*/
				func(buf []byte, updatedAt *time.Time) bool {
					if len(buf) < len(packetDelim) {
						return false
					}
					if bytes.HasSuffix(buf, DATA_PACKET_TAILER) {
						*updatedAt = time.Now()
						return false
					}
					return bytes.HasSuffix(buf, packetDelim)
				},
				2*time.Second,
			)
			recvEndMonotonicTime := timesource.MonotonicNow()
			recvEndTime := g.TimeSource.Now()
			if err != nil {
				g.Logger.Errorf("failed to read data from transport: %v", err)
				cancelFn()
			}
			if timeout {
				g.Logger.Error("timeout when reading data from transport")
				continue
			}

			if len(recvBuf) == 0 {
				continue
			}
			packetBytes := recvBuf[:len(recvBuf)-len(packetDelim)] // without header and tailer

			if err = g.verifyChecksum(packetBytes); err == nil {
				mcuTimestamp := int64(binary.LittleEndian.Uint64(packetBytes[:8]))
				deviceConfig := binary.LittleEndian.Uint32(packetBytes[8 : 8+4])
				gnssEnabled := g.parseGnssAvailibility(deviceConfig)
				channelChunkLength, _, _ := g.getChannelSize(deviceConfig)
				packetInterval := g.parsePacketInterval(deviceConfig)
				sampleRate := g.parseSampleRate(deviceConfig)
				packetLatency := packetInterval + time.Duration(1000/sampleRate*channelChunkLength)*time.Millisecond

				if !g.isDataStreamStable && !gnssEnabled {
					timeDiff := recvEndTime.UnixMilli() - mcuTimestamp - packetLatency.Milliseconds()
					timeDiffSamples = append(timeDiffSamples, timeDiff)
					if len(timeDiffSamples) > STABLE_CHECK_SAMPLES {
						timeDiffSamples = timeDiffSamples[1:]
					}

					if len(timeDiffSamples) == STABLE_CHECK_SAMPLES {
						if minVal, maxVal := lo.Min(timeDiffSamples), lo.Max(timeDiffSamples); math.Abs(float64(maxVal-minVal)) < 5 {
							g.isDataStreamStable = true
							g.fifoBuffer.Reset()
							g.Logger.Infof("data time series stabilized, final time difference = %d ms", timeDiff)
						} else {
							g.Logger.Warnf("waiting for data time series to settle down, this may take a while, current time difference = %d ms", timeDiff)
						}
					} else {
						g.Logger.Warnln("collecting data time series, this may take a while")
					}
				} else if gnssEnabled {
					g.isDataStreamStable = true
				}

				if g.variableAllSet {
					if gnssEnabled && !timeSourceInitialized {
						g.TimeSource.Update(recvEndMonotonicTime, time.UnixMilli(mcuTimestamp).Add(packetLatency), timesource.MonotonicNow)

						g.isDataStreamStable = false
						timeSourceInitialized = true
						g.resetFlags()

						g.Logger.Infof("time synchronized with Explorer built-in GNSS module")
					} else if !timeSourceInitialized {
						g.Logger.Infoln("synchronizing time with NTP servers, it may take a while")
						offset, err := ntpClient.QueryAverage(NTP_MEASUREMENT_ATTEMPTS)
						if err != nil {
							g.Logger.Errorf("failed to synchronize time with NTP server: %v", err)
							if atomic.LoadInt32(&initFlag) == 0 {
								cancelFn()
							} else {
								continue
							}
						} else {
							g.Logger.Infof("time synchronized with NTP server, local monotonic time offset: %d ms", offset.Milliseconds())
						}

						currentMonotonicTime := timesource.MonotonicNow()
						g.TimeSource.Update(currentMonotonicTime, currentMonotonicTime.Add(offset), timesource.MonotonicNow)
						g.isDataStreamStable = false
						timeSourceInitialized = true
						g.resetFlags()
					}

					if atomic.LoadInt32(&initFlag) == 0 {
						atomic.StoreInt32(&initFlag, 1)
						close(readyChan)
						g.deviceStatus.SetStartedAt(g.TimeSource.Now())
					}
				}

				// Handle MCU time jumps (usually caused by Explorer power loss or PC hibernation)
				// 5000 ms is a threshold determined by max packet interval with a safety margin (see getPacketInterval function)
				if (mcuTimestamp < g.prevMcuTimestamp || math.Abs(float64(mcuTimestamp-g.prevMcuTimestamp)) >= 5000) && g.prevMcuTimestamp != 0 {
					g.fifoBuffer.Reset()
					g.resetVariables()
					g.resetFlags()
					g.prevMcuTimestamp = 0
					g.isDataStreamStable = false
					timeDiffSamples = make([]int64, 0, STABLE_CHECK_SAMPLES)
				} else {
					if gnssEnabled && g.isDataStreamStable && g.variableAllSet {
						select {
						case g.timeCalibrationChan <- [2]time.Time{recvEndMonotonicTime, time.UnixMilli(mcuTimestamp).Add(packetLatency)}:
						default:
						}
					}
					g.prevMcuTimestamp = mcuTimestamp
				}

				if g.isDataStreamStable {
					_, _ = g.fifoBuffer.Write(&explorerProtocolPacketV3{
						dataBytes: packetBytes,
						dataTime:  recvEndTime.Add(-packetLatency),
					})
				}
			} else if g.isDataStreamStable && g.variableAllSet {
				g.Logger.Errorln(err)
				g.deviceStatus.IncrementErrors()
			}
		}
	}()

	go func(decodeInterval time.Duration) {
		for timer := time.NewTimer(decodeInterval); ; {
			timer.Reset(decodeInterval)

			select {
			case <-timer.C:
				dataPackets, err := g.fifoBuffer.Read(1)
				if err != nil {
					continue
				}
				dataPacketObj := dataPackets[0]

				mcuTimestamp := int64(binary.LittleEndian.Uint64(dataPacketObj.dataBytes[:8]))
				deviceConfig := binary.LittleEndian.Uint32(dataPacketObj.dataBytes[8 : 8+4])
				g.deviceConfig.SetPacketInterval(g.parsePacketInterval(deviceConfig))
				g.deviceConfig.SetSampleRate(g.parseSampleRate(deviceConfig))
				g.deviceConfig.SetGnssAvailability(g.parseGnssAvailibility(deviceConfig))

				var variableBytes [4]byte
				copy(variableBytes[:], dataPacketObj.dataBytes[8+4:8+4+4])
				g.getVariableData(mcuTimestamp, deviceConfig, variableBytes)
				if !g.variableAllSet {
					g.Logger.Warnln("waiting for device config to be fully collected, this may take a while")
					continue
				}

				if atomic.LoadInt32(&initFlag) == 0 {
					g.Logger.Warnln("waiting for time to be synchronized, this may take a while")
					continue
				}

				gnssEnabled := g.parseGnssAvailibility(deviceConfig)
				timeObj := lo.Ternary(gnssEnabled, time.UnixMilli(mcuTimestamp), dataPacketObj.dataTime)

				g.deviceStatus.IncrementFrames()
				g.deviceStatus.SetUpdatedAt(timeObj)

				channelChunkLength, _, channelData := g.getChannelSize(deviceConfig)
				g.getChannelData(channelData, dataPacketObj.dataBytes[8+4+4:len(dataPacketObj.dataBytes)-1], channelChunkLength)
				g.flagMutex.Lock()
				g.collectedSamples += channelChunkLength
				g.flagMutex.Unlock()

				if g.packetTimeObj.IsZero() {
					g.packetTimeObj = timeObj
				}

				channelCodes := make([]string, len(channelData))
				if len(g.channelDataBuf) != len(channelData) {
					g.channelDataBuf = make([]ChannelData, len(channelData))
				}
				for idx, ch := range channelData {
					chVal := *ch
					channelCodes[idx] = ch.ChannelCode
					g.channelDataBuf[idx].ByteSize = ch.ByteSize
					g.channelDataBuf[idx].ChannelCode = ch.ChannelCode
					g.channelDataBuf[idx].ChannelId = ch.ChannelId
					g.channelDataBuf[idx].DataType = ch.DataType
					g.channelDataBuf[idx].Data = append(g.channelDataBuf[idx].Data, chVal.Data...)
				}
				g.deviceConfig.SetChannelCodes(channelCodes)

				sampleRate := g.deviceConfig.GetSampleRate()
				g.messageBusRealtime.Publish(timeObj, &g.deviceConfig, &g.deviceVariable, g.channelDataBuf)

				g.flagMutex.Lock()
				collectedSamples := g.collectedSamples
				g.flagMutex.Unlock()

				if collectedSamples < sampleRate {
					continue
				} else if collectedSamples == sampleRate {
					g.messageBus.Publish(g.packetTimeObj, &g.deviceConfig, &g.deviceVariable, g.channelDataBuf)
					g.deviceStatus.IncrementMessages()
				} else {
					g.Logger.Warn("collected samples exceeded the sample rate, resetting counters")
					g.resetVariables()
				}

				g.resetFlags()
			case <-subCtx.Done():
				g.Logger.Info("exiting from data packet decoder")
				timer.Stop()
				return
			}
		}
	}(10 * time.Millisecond)

	go func(resyncInterval time.Duration) {
		<-readyChan

		var prevCalibTime time.Time
		for timer := time.NewTimer(resyncInterval); ; {
			select {
			case calibTimeData := <-g.timeCalibrationChan:
				if prevCalibTime.Unix() == calibTimeData[1].Unix() {
					continue
				}
				prevCalibTime = calibTimeData[1]
				g.TimeSource.Update(calibTimeData[0], calibTimeData[1], nil)
			case <-timer.C:
				if deviceConfig := g.GetConfig(); deviceConfig.GetGnssAvailability() || !g.variableAllSet {
					timer.Reset(resyncInterval)
					continue
				}
				g.Logger.Info("re-synchronizing time with NTP servers")
				offset, err := ntpClient.QueryAverage(NTP_MEASUREMENT_ATTEMPTS)
				if err != nil {
					g.Logger.Warnf("error occurred while re-synchronizing time with NTP: %v", err)
					timer.Reset(resyncInterval)
					continue
				}
				timer.Reset(resyncInterval)
				currentMonotonicTime := timesource.MonotonicNow()
				g.TimeSource.Update(currentMonotonicTime, currentMonotonicTime.Add(offset), nil)
				g.Logger.Infof("time synchronized with NTP server, local monotonic time offset: %d ms", offset.Milliseconds())
			case <-subCtx.Done():
				timer.Stop()
				return
			}
		}
	}(NTP_RESYNC_INTERVAL)

	<-readyChan
	return subCtx, cancelFn, nil
}

func (g *ExplorerProtoImplV3) Close() error {
	if g.Transport == nil {
		return errors.New("transport is not opened")
	}

	return g.Transport.Close()
}

func (g *ExplorerProtoImplV3) Subscribe(clientId string, handler EventHandler) error {
	return g.messageBus.Subscribe(clientId, handler)
}

func (g *ExplorerProtoImplV3) Unsubscribe(clientId string) error {
	return g.messageBus.Unsubscribe(clientId)
}

func (g *ExplorerProtoImplV3) SubscribeRealtime(clientId string, handler EventHandler) error {
	return g.messageBusRealtime.Subscribe(clientId, handler)
}

func (g *ExplorerProtoImplV3) UnsubscribeRealtime(clientId string) error {
	return g.messageBusRealtime.Unsubscribe(clientId)
}

func (g *ExplorerProtoImplV3) GetConfig() DeviceConfig {
	return DeviceConfig{
		packetInterval: g.deviceConfig.GetPacketInterval(),
		sampleRate:     g.deviceConfig.GetSampleRate(),
		gnssEnabled:    g.deviceConfig.GetGnssAvailability(),
		channelCodes:   g.deviceConfig.GetChannelCodes(),
		model:          g.deviceConfig.GetModel(),
		protocol:       g.deviceConfig.GetProtocol(),
	}
}

func (g *ExplorerProtoImplV3) GetStatus() DeviceStatus {
	return DeviceStatus{
		startedAt: g.deviceStatus.GetStartedAt(),
		updatedAt: g.deviceStatus.GetUpdatedAt(),
		frames:    g.deviceStatus.GetFrames(),
		errors:    g.deviceStatus.GetErrors(),
		messages:  g.deviceStatus.GetMessages(),
	}
}

func (g *ExplorerProtoImplV3) GetCoordinates(fuzzy bool) (float64, float64, float64, error) {
	lat, err := g.deviceVariable.GetLatitude(fuzzy)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get latitude: %w", err)
	}

	lon, err := g.deviceVariable.GetLongitude(fuzzy)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get longitude: %w", err)
	}

	elv, err := g.deviceVariable.GetElevation()
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get altitude: %w", err)
	}

	return lat, lon, elv, nil
}

func (g *ExplorerProtoImplV3) GetTemperature() (float64, error) {
	temp, err := g.deviceVariable.GetTemperature()
	if err != nil {
		return 0, fmt.Errorf("failed to get temperature: %w", err)
	}
	return temp, nil
}

func (g *ExplorerProtoImplV3) GetDeviceId() string {
	devId, err := g.deviceVariable.GetDeviceId()
	if err != nil {
		return "N/A"
	}
	return fmt.Sprintf("%08X", devId)
}

func (g *ExplorerProtoImplV3) Flush() error {
	return g.Transport.Flush()
}

func (g *ExplorerProtoImplV3) GetMetadata(stationAffiliation, stationDescription, stationCountry, stationPlace, networkCode, stationCode, locationCode string, fuzzyLocation bool) (*metadata.Render, error) {
	latitude, err := g.deviceVariable.GetLatitude(fuzzyLocation)
	if err != nil {
		return nil, fmt.Errorf("failed to get latitude: %w", err)
	}
	longitude, err := g.deviceVariable.GetLongitude(fuzzyLocation)
	if err != nil {
		return nil, fmt.Errorf("failed to get longitude: %w", err)
	}
	elevation, err := g.deviceVariable.GetElevation()
	if err != nil {
		return nil, fmt.Errorf("failed to get altitude: %w", err)
	}
	return metadata.New(g.ExplorerOptions.Model, metadata.Options{
		ChannelCodes:       g.deviceConfig.GetChannelCodes(),
		StartTime:          g.deviceStatus.GetStartedAt(),
		SampleRate:         g.deviceConfig.GetSampleRate(),
		Latitude:           latitude,
		Longitude:          longitude,
		Elevation:          elevation,
		NetworkCode:        networkCode,
		StationCode:        stationCode,
		LocationCode:       locationCode,
		StationAffiliation: stationAffiliation,
		StationDescription: stationDescription,
		StationCountry:     stationCountry,
		StationPlace:       stationPlace,
	})
}
