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

type ExplorerProtoImplV2 struct {
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

	needUpdateTimeSource bool
	variableAllSet       bool

	deviceStatus   DeviceStatus
	deviceConfig   DeviceConfig
	deviceVariable DeviceVariable
	channelDataBuf []ChannelData
}

func (g *ExplorerProtoImplV2) resetVariables() {
	g.needUpdateTimeSource = true
	g.variableAllSet = false
	g.deviceVariable.Reset()
}

func (g *ExplorerProtoImplV2) getPacketSize(headerSize, channelSize int) int {
	return headerSize + // header
		int(unsafe.Sizeof(int64(0))+ // timestamp
			unsafe.Sizeof(uint32(0))+ // variable data
			uintptr(channelSize)*unsafe.Sizeof(int32(0))+ // channel 1
			uintptr(channelSize)*unsafe.Sizeof(int32(0))+ // channel 2
			uintptr(channelSize)*unsafe.Sizeof(int32(0))+ // channel 3
			unsafe.Sizeof(uint8(0))) // checksum
}

func (g *ExplorerProtoImplV2) getTimestamp(mcuTimestamp int64) int64 {
	if g.deviceConfig.GetGnssAvailability() {
		return mcuTimestamp
	}

	g.timeDiffMutex.Lock()
	timestamp := mcuTimestamp + g.timeDiff4NonGnssMode
	g.timeDiffMutex.Unlock()

	return timestamp
}

func (g *ExplorerProtoImplV2) getVariableData(mcuTimestamp int64, variableBytes uint32) {
	gnssEnabled := g.deviceConfig.GetGnssAvailability()
	switch (mcuTimestamp / 1000) % 4 {
	case 0:
		deviceId := variableBytes & 0x7FFFFFFF
		g.deviceVariable.SetDeviceId(&deviceId)
		gnssEnable := (variableBytes&0x80000000 != 0)
		if !gnssEnable {
			g.variableAllSet = true
		}
		g.deviceConfig.SetGnssAvailability(gnssEnable)
	case 1:
		if gnssEnabled {
			n := float64(math.Float32frombits(variableBytes))
			g.deviceVariable.SetLatitude(&n)
		} else {
			g.deviceVariable.SetLatitude(&g.ExplorerOptions.Latitude)
		}
	case 2:
		if gnssEnabled {
			n := float64(math.Float32frombits(variableBytes))
			g.deviceVariable.SetLongitude(&n)
		} else {
			g.deviceVariable.SetLongitude(&g.ExplorerOptions.Longitude)
		}
	case 3:
		if gnssEnabled {
			n := float64(math.Float32frombits(variableBytes))
			g.deviceVariable.SetElevation(&n)
		} else {
			g.deviceVariable.SetElevation(&g.ExplorerOptions.Elevation)
		}
	}

	if _, err := g.deviceVariable.GetLatitude(false); err != nil {
		return
	}
	if _, err := g.deviceVariable.GetLongitude(false); err != nil {
		return
	}
	if _, err := g.deviceVariable.GetElevation(); err != nil {
		return
	}
	if gnssEnabled {
		g.variableAllSet = true
	}
}

func (g *ExplorerProtoImplV2) getChannelData(packetBytes []byte, headerSize, channelSize int) error {
	if len(g.channelDataBuf) != 3 {
		g.channelDataBuf = make([]ChannelData, 3)
	}

	zOffset := headerSize + int(unsafe.Sizeof(int64(0))) + int(unsafe.Sizeof(uint32(0)))
	zAxisData := make([]int32, channelSize)
	eOffset := zOffset + (channelSize)*int(unsafe.Sizeof(int32(0)))
	eAxisData := make([]int32, channelSize)
	nOffset := eOffset + (channelSize)*int(unsafe.Sizeof(int32(0)))
	nAxisData := make([]int32, channelSize)

	err := binary.Read(bytes.NewReader(packetBytes[zOffset:eOffset]), binary.LittleEndian, &zAxisData)
	if err != nil {
		return fmt.Errorf("failed to read z-axis data: %w", err)
	}
	err = binary.Read(bytes.NewReader(packetBytes[eOffset:nOffset]), binary.LittleEndian, &eAxisData)
	if err != nil {
		return fmt.Errorf("failed to read e-axis data: %w", err)
	}
	err = binary.Read(bytes.NewReader(packetBytes[nOffset:len(packetBytes)-1]), binary.LittleEndian, &nAxisData)
	if err != nil {
		return fmt.Errorf("failed to read n-axis data: %w", err)
	}

	for i := 0; i < len(g.channelDataBuf); i++ {
		channelId := i + 1
		g.channelDataBuf[i].ChannelCode = fmt.Sprintf("CH%d", channelId)
		if i < len(g.ChannelCodes) {
			g.channelDataBuf[i].ChannelCode = g.ChannelCodes[i]
		}
		g.channelDataBuf[i].ChannelId = channelId
		g.channelDataBuf[i].ByteSize = 4
		g.channelDataBuf[i].DataType = "int32"
		for j := 0; j < channelSize; j++ {
			switch i {
			case 0:
				g.channelDataBuf[i].Data = append(g.channelDataBuf[i].Data, zAxisData[j])
			case 1:
				g.channelDataBuf[i].Data = append(g.channelDataBuf[i].Data, eAxisData[j])
			case 2:
				g.channelDataBuf[i].Data = append(g.channelDataBuf[i].Data, nAxisData[j])
			}
		}
	}

	var currentChannelCodes []string
	for _, channelData := range g.channelDataBuf {
		currentChannelCodes = append(currentChannelCodes, channelData.ChannelCode)
	}
	g.deviceConfig.SetChannelCodes(currentChannelCodes)

	return nil
}

func (g *ExplorerProtoImplV2) verifyChecksum(packetData, header []byte) error {
	if len(packetData) <= len(header) {
		return errors.New("invalid packet length")
	}
	recvChecksum := packetData[len(packetData)-1]
	calcChecksum := uint8(0)
	for _, b := range packetData[len(header) : len(packetData)-1] {
		calcChecksum ^= b
	}
	if recvChecksum != calcChecksum {
		return fmt.Errorf("invalid checksum: expected %v, got %v", recvChecksum, calcChecksum)
	}
	return nil
}

func (g *ExplorerProtoImplV2) Open(ctx context.Context) (context.Context, context.CancelFunc, error) {
	if g.Transport == nil {
		return nil, nil, errors.New("transport is not opened")
	}
	if g.Logger == nil {
		return nil, nil, errors.New("logger is not set")
	}

	if err := g.Transport.Open(); err != nil {
		return nil, nil, fmt.Errorf("failed to open transport: %w", err)
	}
	ntpClient, err := ntpclient.New(g.NtpOptions.Endpoint, g.NtpOptions.Retry, g.NtpOptions.ReadTimeout)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ntp client: %w", err)
	}

	subCtx, cancelFn := context.WithCancel(ctx)

	const (
		// In v2 protocol, each packet contains 3 channels, n samples per channel.
		// The packet is sent at an interval of (1000 / sample rate) milliseconds.
		// Set n = 5 (also in Explorer) fits the common sample rates (25, 50, 100, 125 Hz).
		DATA_PACKET_CHANNEL_SIZE = 5
		ALLOWED_JITTER_MS        = 5
	)

	DATA_PACKET_HEADER := []byte{0xFA, 0xDE}
	packetSize := g.getPacketSize(len(DATA_PACKET_HEADER), DATA_PACKET_CHANNEL_SIZE)
	g.fifoBuffer = fifo.New[byte](10 * packetSize)
	g.messageBus = message.NewBus[EventHandler](EXPLORER_STREAM_TOPIC, 1024)
	g.deviceStatus.SetUpdatedAt(time.Unix(0, 0))
	g.deviceConfig.SetProtocol(g.ExplorerOptions.Protocol)
	g.deviceConfig.SetModel(g.ExplorerOptions.Model)
	g.timeCalibrationChan4GnssMode = make(chan [2]time.Time)

	var initFlag int32
	atomic.StoreInt32(&initFlag, 0)
	readyChan := make(chan struct{})

	go func() {
		timeDiffSamples := make([]int64, 0, STABLE_CHECK_SAMPLES)
		g.isTimeDiff4NonGnssModeStable = false
		g.needUpdateTimeSource = true

		buf := make([]byte, packetSize*2)
		_ = g.Flush()

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
			n, err := g.Transport.Read(buf)
			recvEndTime := g.TimeSource.Now()
			if err != nil {
				g.Logger.Errorf("failed to read data from transport: %v", err)
				cancelFn()
			}
			recvBuf := buf[:n]

			packetLatency := recvEndTime.Sub(recvStartTime)
			if headerIdx := bytes.Index(recvBuf, DATA_PACKET_HEADER); headerIdx != -1 && len(recvBuf) >= headerIdx+packetSize {
				if err = g.verifyChecksum(recvBuf[headerIdx:headerIdx+packetSize], DATA_PACKET_HEADER); err == nil {
					mcuTimestamp := int64(binary.LittleEndian.Uint64(recvBuf[headerIdx+len(DATA_PACKET_HEADER) : headerIdx+len(DATA_PACKET_HEADER)+int(unsafe.Sizeof(int64(0)))]))

					estimatedTransportLatency := g.Transport.GetLatency(len(recvBuf))
					packetLatency += estimatedTransportLatency
					timeDiff := recvEndTime.UnixMilli() - mcuTimestamp - packetLatency.Milliseconds()

					if !g.isTimeDiff4NonGnssModeStable {
						timeDiffSamples = append(timeDiffSamples, timeDiff)
						if len(timeDiffSamples) > STABLE_CHECK_SAMPLES {
							timeDiffSamples = timeDiffSamples[1:]
						}

						if len(timeDiffSamples) == STABLE_CHECK_SAMPLES {
							if minVal, maxVal := lo.Min(timeDiffSamples), lo.Max(timeDiffSamples); math.Abs(float64(maxVal-minVal)) <= 5 {
								g.isTimeDiff4NonGnssModeStable = true
								g.Logger.Infof("data time series stabilized: time difference = %d ms", timeDiff)
							} else {
								g.Logger.Warnln("waiting for data time series to settle down, this may take a while")
							}
						} else {
							if err = g.Flush(); err != nil {
								g.Logger.Errorf("failed to flush transport: %v", err)
								cancelFn()
							}
							g.Logger.Warnln("collecting data time series, this may take a while")
						}
					}

					if g.deviceConfig.GetSampleRate() > 0 && g.variableAllSet {
						if g.deviceConfig.GetGnssAvailability() && g.needUpdateTimeSource {
							g.TimeSource.Update(recvEndTime, time.UnixMilli(mcuTimestamp).Add(packetLatency))

							g.needUpdateTimeSource = false
							g.isTimeDiff4NonGnssModeStable = false

							g.Logger.Infof("time synchronized with Explorer built-in GNSS module")
						} else if g.needUpdateTimeSource {
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
							g.needUpdateTimeSource = false
							g.isTimeDiff4NonGnssModeStable = false
						}
					}

					if atomic.LoadInt32(&initFlag) == 0 {
						atomic.StoreInt32(&initFlag, 1)
						close(readyChan)
						g.deviceStatus.SetStartedAt(g.TimeSource.Now())
					}

					if g.deviceConfig.GetSampleRate() > 0 && g.variableAllSet && g.deviceConfig.GetGnssAvailability() {
						select {
						case g.timeCalibrationChan4GnssMode <- [2]time.Time{recvEndTime, time.UnixMilli(mcuTimestamp).Add(packetLatency)}:
						default:
						}
					}

					// Compensate for oscillator drift on the AnyShake Explorer board (NTP mode only)
					if g.deviceConfig.GetSampleRate() > 0 && g.variableAllSet && !g.deviceConfig.GetGnssAvailability() {
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
					// 5000 ms is a threshold determined by max packet interval with a minimum sample rate of 1 Hz
					if (mcuTimestamp < g.prevMcuTimestamp || math.Abs(float64(mcuTimestamp)-float64(g.prevMcuTimestamp)) >= 5000) && g.prevMcuTimestamp != 0 {
						g.fifoBuffer.Reset()
						g.resetVariables()
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
		var (
			expectedNextMcuTimestamp int64
			collectedTimestampArr    []int64
		)
		for timer := time.NewTimer(decodeInterval); ; {
			timer.Reset(decodeInterval)

			select {
			case <-timer.C:
				dataPacket, err := g.fifoBuffer.Peek(DATA_PACKET_HEADER, packetSize)
				if err != nil {
					continue
				}

				mcuTimestamp := int64(binary.LittleEndian.Uint64(dataPacket[2:10]))
				variableData := binary.LittleEndian.Uint32(dataPacket[10:14])
				g.getVariableData(mcuTimestamp, variableData)

				if !g.variableAllSet {
					g.Logger.Warnln("waiting for device config to be fully collected, this may take a while")
					continue
				}

				if err = g.verifyChecksum(dataPacket, DATA_PACKET_HEADER); err != nil {
					g.Logger.Errorln(err)
					g.deviceStatus.IncrementErrors()
					continue
				}

				timestamp := g.getTimestamp(mcuTimestamp)
				if expectedNextMcuTimestamp == 0 {
					expectedNextMcuTimestamp = mcuTimestamp + 1000
				} else {
					collectedTimestampArr = append(collectedTimestampArr, timestamp)
					err = g.getChannelData(dataPacket, len(DATA_PACKET_HEADER), DATA_PACKET_CHANNEL_SIZE)
					if err != nil {
						g.Logger.Errorf("failed to get channel data: %v", err)
						g.deviceStatus.IncrementErrors()
						continue
					}
				}

				if atomic.LoadInt32(&initFlag) == 0 {
					g.Logger.Warnln("waiting for time to be synchronized, this may take a while")
					continue
				}

				if math.Abs(float64(mcuTimestamp-expectedNextMcuTimestamp)) <= ALLOWED_JITTER_MS {
					// Update the next tick even if the buffer is empty
					expectedNextMcuTimestamp = mcuTimestamp + time.Second.Milliseconds()
					if len(collectedTimestampArr) == 0 {
						continue
					}

					sampleRate := len(collectedTimestampArr) * DATA_PACKET_CHANNEL_SIZE
					g.deviceConfig.SetSampleRate(sampleRate)
					g.deviceConfig.SetPacketInterval(time.Duration((1000/sampleRate)*DATA_PACKET_CHANNEL_SIZE) * time.Millisecond)
					packetTimestamp := collectedTimestampArr[0]
					g.messageBus.Publish(time.UnixMilli(packetTimestamp), &g.deviceConfig, &g.deviceVariable, g.channelDataBuf)
					g.deviceStatus.IncrementMessages()

					collectedTimestampArr = []int64{}
					g.channelDataBuf = []ChannelData{}
				} else if expectedNextMcuTimestamp-mcuTimestamp > time.Second.Milliseconds()+ALLOWED_JITTER_MS || expectedNextMcuTimestamp-mcuTimestamp < 0 {
					g.Logger.Warnf("jitter detected, discarding this packet, expected %v, got %v", g.getTimestamp(expectedNextMcuTimestamp), timestamp)
					g.deviceStatus.IncrementErrors()
					// Update the next tick, clear the buffer if the jitter exceeds the threshold
					expectedNextMcuTimestamp = mcuTimestamp + time.Second.Milliseconds()
					collectedTimestampArr = []int64{}
					g.channelDataBuf = []ChannelData{}
				}

				g.deviceStatus.IncrementFrames()
				g.deviceStatus.SetUpdatedAt(time.UnixMilli(int64(timestamp)))
			case <-subCtx.Done():
				g.Logger.Info("exiting from data packet decoder")
				timer.Stop()
				return
			}
		}
	}(5 * time.Millisecond)

	go func() {
		<-readyChan

		getNextInterval := func() time.Duration {
			now := g.TimeSource.Now()
			interval := time.Hour
			if g.deviceConfig.GetGnssAvailability() {
				interval = time.Minute
			}

			next := now.Truncate(interval).Add(interval)
			if g.deviceConfig.GetGnssAvailability() && next.Hour() == 0 && next.Minute() == 0 {
				next = next.Add(interval)
			}

			return time.Until(next)
		}
		for timer := time.NewTimer(getNextInterval()); ; {
			timer.Reset(getNextInterval())

			select {
			case <-timer.C:
				if g.deviceConfig.GetGnssAvailability() && g.TimeSource.Now().Hour() == 0 {
					continue
				}

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

func (g *ExplorerProtoImplV2) Close() error {
	if g.Transport == nil {
		return errors.New("transport is not opened")
	}

	return g.Transport.Close()
}

func (g *ExplorerProtoImplV2) Subscribe(clientId string, handler EventHandler) error {
	return g.messageBus.Subscribe(clientId, handler)
}

func (g *ExplorerProtoImplV2) Unsubscribe(clientId string) error {
	return g.messageBus.Unsubscribe(clientId)
}

func (g *ExplorerProtoImplV2) GetConfig() DeviceConfig {
	return DeviceConfig{
		packetInterval: g.deviceConfig.GetPacketInterval(),
		sampleRate:     g.deviceConfig.GetSampleRate(),
		gnssEnabled:    g.deviceConfig.GetGnssAvailability(),
		channelCodes:   g.deviceConfig.GetChannelCodes(),
		model:          g.deviceConfig.GetModel(),
		protocol:       g.deviceConfig.GetProtocol(),
	}
}

func (g *ExplorerProtoImplV2) GetStatus() DeviceStatus {
	return DeviceStatus{
		startedAt: g.deviceStatus.GetStartedAt(),
		updatedAt: g.deviceStatus.GetUpdatedAt(),
		frames:    g.deviceStatus.GetFrames(),
		errors:    g.deviceStatus.GetErrors(),
		messages:  g.deviceStatus.GetMessages(),
	}
}

func (g *ExplorerProtoImplV2) GetCoordinates(fuzzy bool) (float64, float64, float64, error) {
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

func (g *ExplorerProtoImplV2) GetTemperature() (float64, error) {
	temp, err := g.deviceVariable.GetTemperature()
	if err != nil {
		return 0, fmt.Errorf("failed to get temperature: %w", err)
	}
	return temp, nil
}

func (g *ExplorerProtoImplV2) GetDeviceId() string {
	devId, err := g.deviceVariable.GetDeviceId()
	if err != nil {
		return "N/A"
	}
	return fmt.Sprintf("%08X", devId)
}

func (g *ExplorerProtoImplV2) Flush() error {
	return g.Transport.Flush()
}

func (g *ExplorerProtoImplV2) GetMetadata(stationAffiliation, stationDescription, stationCountry, stationPlace, networkCode, stationCode, locationCode string, fuzzyLocation bool) (metadata.IMetadata, error) {
	latitude, longitude, elevation, err := g.GetCoordinates(fuzzyLocation)
	if err != nil {
		return nil, err
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
