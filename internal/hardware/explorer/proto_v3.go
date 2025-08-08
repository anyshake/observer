package explorer

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/anyshake/observer/internal/hardware/explorer/metadata"
	"github.com/anyshake/observer/pkg/fifo"
	"github.com/anyshake/observer/pkg/message"
	"github.com/anyshake/observer/pkg/ntpclient"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/anyshake/observer/pkg/transport"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type ExplorerProtoImplV3 struct {
	ChannelCodes    []string
	ExplorerOptions ExplorerOptions
	NtpOptions      NtpOptions
	Logger          *logrus.Entry
	TimeSource      *timesource.Source

	Transport  transport.ITransport
	fifoBuffer fifo.Buffer[byte]
	messageBus message.Bus[EventHandler]

	timeDiffMutex                sync.Mutex
	prevMcuTimestamp             int64
	prevTimestamp4NonGnssMode    int64
	timeDiff4NonGnssMode         int64
	prevTimeOffset4NonGnssMode   *int64
	isTimeDiff4NonGnssModeStable bool
	timeCalibrationChan4GnssMode chan [2]time.Time

	flagMutex            sync.Mutex
	needUpdateTimeSource bool
	variableAllSet       bool
	collectedSamples     int
	packetTimeObj        time.Time

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
	g.needUpdateTimeSource = true
	g.variableAllSet = false
	g.flagMutex.Unlock()

	g.deviceVariable.Reset()
}

func (g *ExplorerProtoImplV3) getTimestamp(mcuTimestamp int64) int64 {
	if g.deviceConfig.GetGnssAvailability() {
		return mcuTimestamp
	}

	g.timeDiffMutex.Lock()
	timestamp := mcuTimestamp + g.timeDiff4NonGnssMode
	g.timeDiffMutex.Unlock()

	return timestamp
}

func (g *ExplorerProtoImplV3) getPacketInterval(deviceConfig uint32) {
	DATA_PACKET_PACKET_INTERVAL := []int{100, 200, 500, 1000}
	g.deviceConfig.SetPacketInterval(time.Duration(DATA_PACKET_PACKET_INTERVAL[(deviceConfig>>30)&0x3]) * time.Millisecond)
}

func (g *ExplorerProtoImplV3) getSampleRate(deviceConfig uint32) {
	DATA_PACKET_SAMPLE_RATES := []int{10, 20, 50, 100, 200, 250, 500, 1000}
	g.deviceConfig.SetSampleRate(DATA_PACKET_SAMPLE_RATES[(deviceConfig>>27)&0x7])
}

func (g *ExplorerProtoImplV3) getGnssAvailibility(deviceConfig uint32) {
	g.deviceConfig.SetGnssAvailability(((deviceConfig >> 26) & 0x1) == 1)
}

func (g *ExplorerProtoImplV3) getChannelSize(deviceConfig uint32) (channelChunkLength, totalChannelSize int, channelData []*ChannelData) {
	DATA_PACKET_CHANNEL_TYPE := []string{"disabled", "int16", "int24", "int32"}
	channelChunkLength = int(g.deviceConfig.GetPacketInterval().Milliseconds()) / (1000 / g.deviceConfig.GetSampleRate())

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

func (g *ExplorerProtoImplV3) verifyChecksum(packetData, header, tailer []byte) error {
	if len(packetData) <= len(header)+len(tailer) {
		return errors.New("invalid packet length")
	}
	recvChecksum := packetData[len(packetData)-len(tailer)-1]
	calcChecksum := uint8(0)
	for _, b := range packetData[len(header) : len(packetData)-len(tailer)-1] {
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

	if g.Logger == nil {
		return nil, nil, errors.New("logger is not set")
	}
	ntpClient, err := ntpclient.New(g.NtpOptions.Endpoint, g.NtpOptions.Retry, g.NtpOptions.ReadTimeout)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ntp client: %w", err)
	}

	subCtx, cancelFn := context.WithCancel(ctx)

	// Assume that the longest packet interval is 1000 ms
	// With 8 channels and 1000 samples per second per channel in int32
	// That would be 8 channels * 4 bytes * (1000 ms / (1000 / 1000 SPS))
	// Set to 655350 would be safe enough to avoid buffer overflows
	g.fifoBuffer = fifo.New[byte](655350)
	g.messageBus = message.NewBus[EventHandler](EXPLORER_STREAM_TOPIC, 1024)
	g.deviceStatus.SetUpdatedAt(time.Unix(0, 0))
	g.deviceConfig.SetProtocol(g.ExplorerOptions.Protocol)
	g.deviceConfig.SetModel(g.ExplorerOptions.Model)
	g.timeCalibrationChan4GnssMode = make(chan [2]time.Time)

	var initFlag int32
	atomic.StoreInt32(&initFlag, 0)
	readyChan := make(chan struct{})

	DATA_PACKET_HEADER := []byte{0x01, 0xFE}
	DATA_PACKET_TAILER := []byte{0xEF, 0x10}

	go func() {
		timeDiffSamples := make([]int64, 0, STABLE_CHECK_SAMPLES)
		g.isTimeDiff4NonGnssModeStable = false

		g.flagMutex.Lock()
		g.needUpdateTimeSource = true
		g.flagMutex.Unlock()

		for {
			select {
			case <-subCtx.Done():
				g.Logger.Info("exiting from data packet reader")
				if atomic.LoadInt32(&initFlag) == 0 {
					close(readyChan)
				}
				return
			default:
			}

			recvStartTime := g.TimeSource.Now()
			recvBuf, timeout, recvElapsed, err := g.Transport.ReadUntil(subCtx, DATA_PACKET_TAILER, 32000, 5*time.Second)
			recvEndTime := g.TimeSource.Now()
			if err != nil {
				g.Logger.Errorf("failed to read data from transport: %v", err)
				cancelFn()
			}
			if timeout {
				g.Logger.Error("timeout when reading data from transport")
				continue
			}

			totalLatency := recvEndTime.Sub(recvStartTime)
			if headerIdx := bytes.Index(recvBuf, DATA_PACKET_HEADER); headerIdx != -1 {
				if err = g.verifyChecksum(recvBuf[headerIdx:], DATA_PACKET_HEADER, DATA_PACKET_TAILER); err == nil && len(recvBuf) > headerIdx+len(DATA_PACKET_HEADER)+int(unsafe.Sizeof(int64(0))) {
					mcuTimestamp := int64(binary.LittleEndian.Uint64(recvBuf[headerIdx+len(DATA_PACKET_HEADER) : headerIdx+len(DATA_PACKET_HEADER)+int(unsafe.Sizeof(int64(0)))]))

					estimatedTransportLatency := g.Transport.GetLatency(len(recvBuf))
					totalLatency += lo.Ternary(estimatedTransportLatency > recvElapsed, estimatedTransportLatency, recvElapsed)
					timeDiff := recvEndTime.UnixMilli() - mcuTimestamp - totalLatency.Milliseconds()

					if !g.isTimeDiff4NonGnssModeStable {
						timeDiffSamples = append(timeDiffSamples, timeDiff)
						if len(timeDiffSamples) > STABLE_CHECK_SAMPLES {
							timeDiffSamples = timeDiffSamples[1:]
						}

						if len(timeDiffSamples) == STABLE_CHECK_SAMPLES {
							if minVal, maxVal := lo.Min(timeDiffSamples), lo.Max(timeDiffSamples); math.Abs(float64(maxVal-minVal)) <= 10 {
								g.isTimeDiff4NonGnssModeStable = true
								g.Logger.Infof("data time series stabilized: time difference = %d ms", timeDiff)
							} else {
								g.Logger.Warnln("waiting for data time series to settle down, this may take a while")
							}
						} else {
							if err = g.Flush(); err != nil {
								g.Logger.Errorf("failed to flush transport: %v", err)
								cancelFn()
							} else {
								g.Logger.Warnln("collecting data time series, this may take a while")
							}
						}
					}

					if g.deviceConfig.GetSampleRate() > 0 && g.variableAllSet {
						g.flagMutex.Lock()
						needUpdateTimeSource := g.needUpdateTimeSource
						g.flagMutex.Unlock()

						if g.deviceConfig.GetGnssAvailability() && needUpdateTimeSource {
							g.TimeSource.Update(recvEndTime, time.UnixMilli(mcuTimestamp).Add(totalLatency))

							g.isTimeDiff4NonGnssModeStable = false
							g.resetFlags()

							g.flagMutex.Lock()
							g.needUpdateTimeSource = false
							g.flagMutex.Unlock()

							g.Logger.Infof("time synchronized with Explorer built-in GNSS module")
						} else if needUpdateTimeSource {
							g.Logger.Infof("synchronizing time with NTP server: %s", g.NtpOptions.Endpoint)
							res, err := ntpClient.Query()
							if err != nil {
								g.Logger.Errorf("failed to synchronize time with NTP server: %v", err)
								if atomic.LoadInt32(&initFlag) == 0 {
									cancelFn()
								} else {
									continue
								}
							} else {
								g.Logger.Infof("time synchronized with NTP server, local time offset: %d ms", res.ClockOffset.Milliseconds())
							}

							currentTime := time.Now()
							g.TimeSource.Update(currentTime, currentTime.Add(res.ClockOffset))
							g.isTimeDiff4NonGnssModeStable = false
							g.resetFlags()

							g.flagMutex.Lock()
							g.needUpdateTimeSource = false
							g.flagMutex.Unlock()
						}

						if atomic.LoadInt32(&initFlag) == 0 {
							atomic.StoreInt32(&initFlag, 1)
							close(readyChan)
							g.deviceStatus.SetStartedAt(g.TimeSource.Now())
						}

						if g.deviceConfig.GetGnssAvailability() {
							select {
							case g.timeCalibrationChan4GnssMode <- [2]time.Time{recvEndTime, time.UnixMilli(mcuTimestamp).Add(totalLatency)}:
							default:
							}
						}
					}

					// Compensate for oscillator drift on the AnyShake Explorer board (NTP mode only)
					if g.deviceConfig.GetSampleRate() > 0 && !g.deviceConfig.GetGnssAvailability() {
						timeOffset := g.getTimestamp(mcuTimestamp) - g.TimeSource.Now().UnixMilli()
						if g.prevTimeOffset4NonGnssMode == nil {
							g.prevTimeOffset4NonGnssMode = &timeOffset
						}
						if math.Abs(float64(timeOffset)-float64(*g.prevTimeOffset4NonGnssMode)) > 1 {
							g.timeDiff4NonGnssMode = timeDiff
							g.prevTimeOffset4NonGnssMode = &timeOffset
						}
					}

					g.timeDiffMutex.Lock()

					if g.timeDiff4NonGnssMode == 0 && timeDiff != 0 {
						g.timeDiff4NonGnssMode = timeDiff
					}

					// Handle MCU time jumps (usually caused by Explorer power loss or PC hibernation)
					// 1500 ms is a threshold determined by max packet interval with a safety margin (see getPacketInterval function)
					if (mcuTimestamp < g.prevMcuTimestamp || math.Abs(float64(mcuTimestamp)-float64(g.prevMcuTimestamp)) >= 1500) && g.prevMcuTimestamp != 0 {
						g.fifoBuffer.Reset()
						g.resetVariables()
						g.resetFlags()
						g.timeDiff4NonGnssMode = 0
						g.isTimeDiff4NonGnssModeStable = false
						timeDiffSamples = make([]int64, 0, STABLE_CHECK_SAMPLES)
					}

					g.prevMcuTimestamp = mcuTimestamp
					g.prevTimestamp4NonGnssMode = g.prevMcuTimestamp + g.timeDiff4NonGnssMode

					g.timeDiffMutex.Unlock()
				}
			}

			if g.isTimeDiff4NonGnssModeStable {
				_, _ = g.fifoBuffer.Write(recvBuf...)
			}
		}
	}()

	go func(decodeInterval time.Duration) {
		for timer := time.NewTimer(decodeInterval); ; {
			timer.Reset(decodeInterval)

			select {
			case <-timer.C:
				// Fixed 18 bytes for header + timestamp + device config + variable data
				packetFixedSection, err := g.fifoBuffer.Peek(DATA_PACKET_HEADER, len(DATA_PACKET_HEADER)+16)
				if err != nil {
					continue
				}

				mcuTimestamp := int64(binary.LittleEndian.Uint64(packetFixedSection[2:10]))
				deviceConfig := binary.LittleEndian.Uint32(packetFixedSection[10:14])
				g.getPacketInterval(deviceConfig)
				g.getSampleRate(deviceConfig)
				g.getGnssAvailibility(deviceConfig)

				timestamp := g.getTimestamp(mcuTimestamp)
				timeObj := time.UnixMilli(int64(timestamp))

				var variableBytes [4]byte
				copy(variableBytes[:], packetFixedSection[14:18])
				g.getVariableData(mcuTimestamp, deviceConfig, variableBytes)

				// Calculate channel data size and read data remaining (channel data + checksum + tailer)
				channelChunkLength, channelSize, channelData := g.getChannelSize(deviceConfig)
				readSize := channelSize + 1 + len(DATA_PACKET_TAILER)
				for g.fifoBuffer.Len() < readSize {
					time.Sleep(10 * time.Millisecond)
				}
				channelDataSection, err := g.fifoBuffer.Read(readSize)
				if err != nil {
					g.Logger.Errorf("failed to read channel data: %v", err)
					g.deviceStatus.IncrementErrors()
					continue
				}

				recvTailer := channelDataSection[len(channelDataSection)-2:]
				if !bytes.Equal(recvTailer, DATA_PACKET_TAILER) {
					g.Logger.Errorf("tailer mismatch, expected %v, got %v", DATA_PACKET_TAILER, channelDataSection[len(channelDataSection)-2:])
					g.deviceStatus.IncrementErrors()
					continue
				}

				if err = g.verifyChecksum(
					bytes.Join([][]byte{packetFixedSection, channelDataSection}, nil),
					DATA_PACKET_HEADER, DATA_PACKET_TAILER,
				); err != nil {
					g.Logger.Errorln(err)
					g.deviceStatus.IncrementErrors()
					continue
				}

				g.deviceStatus.IncrementFrames()
				g.deviceStatus.SetUpdatedAt(timeObj)

				if !g.variableAllSet {
					g.Logger.Warnln("waiting for device config to be fully collected, this may take a while")
					continue
				}

				if atomic.LoadInt32(&initFlag) == 0 {
					g.Logger.Warnln("waiting for time to be synchronized, this may take a while")
					continue
				}

				g.getChannelData(channelData, channelDataSection[:len(channelDataSection)-1-len(recvTailer)], channelChunkLength)
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

	go func() {
		<-readyChan

		getNextHour := func() time.Duration {
			now := g.TimeSource.Now()
			nextHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, time.UTC)
			return time.Until(nextHour)
		}
		for timer := time.NewTimer(getNextHour()); ; {
			timer.Reset(getNextHour())

			select {
			case <-timer.C:
				if g.deviceConfig.GetGnssAvailability() {
					select {
					case calibTimeData := <-g.timeCalibrationChan4GnssMode:
						g.TimeSource.Update(calibTimeData[0], calibTimeData[1])
					case <-time.After(time.Second):
						g.Logger.Warn("no GNSS calibration timestamp received within 1 second, skipping")
					}
				} else {
					res, err := ntpClient.Query()
					if err != nil {
						g.Logger.Warnf("error occurred while re-synchronizing time: %v", err)
						continue
					}
					currentTime := time.Now()
					g.TimeSource.Update(currentTime, currentTime.Add(res.ClockOffset))
				}
			case <-subCtx.Done():
				timer.Stop()
				return
			}
		}
	}()

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

func (g *ExplorerProtoImplV3) GetMetadata(stationAffiliation, stationDescription, stationCountry, stationPlace, networkCode, stationCode, locationCode string, fuzzyLocation bool) (metadata.IMetadata, error) {
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
	return metadata.New(g.deviceConfig.GetModel(), metadata.Options{
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
