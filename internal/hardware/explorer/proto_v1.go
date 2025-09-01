package explorer

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"
	"unsafe"

	"github.com/anyshake/observer/internal/hardware/explorer/metadata"
	"github.com/anyshake/observer/pkg/fifo"
	"github.com/anyshake/observer/pkg/message"
	"github.com/anyshake/observer/pkg/ntpclient"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/anyshake/observer/pkg/transport"
	"github.com/sirupsen/logrus"
)

type ExplorerProtoImplV1 struct {
	ChannelCodes    []string
	ExplorerOptions ExplorerOptions
	NtpOptions      NtpOptions
	Logger          *logrus.Entry
	TimeSource      *timesource.Source

	Transport  transport.ITransport
	fifoBuffer fifo.Buffer[byte]
	messageBus message.Bus[EventHandler]

	deviceStatus   DeviceStatus
	deviceConfig   DeviceConfig
	deviceVariable DeviceVariable
	channelDataBuf []ChannelData
}

func (g *ExplorerProtoImplV1) getPacketSize(headerSize, channelSize int) int {
	return headerSize + // header
		int(uintptr(channelSize)*unsafe.Sizeof(int32(0))+ // channel 1
			uintptr(channelSize)*unsafe.Sizeof(int32(0))+ // channel 2
			uintptr(channelSize)*unsafe.Sizeof(int32(0))) + // channel 3
		int(3*unsafe.Sizeof(uint8(0))) + // checksum
		1 // padding
}

func (g *ExplorerProtoImplV1) fixSampleRate(currentSampleRate int) int {
	targetSampleRates := []int{
		5, 10, 25, 50, 75, 100,
		125, 150, 175, 200, 225,
		250, 275, 300, 325, 350,
		375, 400, 425, 450, 475,
		500, 525, 550, 575, 600,
		625, 650, 675, 700, 725,
		750, 775, 800, 825, 850,
		875, 900, 925, 950, 975,
		1000,
	}

	closest := targetSampleRates[0]
	minDiff := math.Abs(float64(currentSampleRate - closest))

	for _, target := range targetSampleRates {
		diff := math.Abs(float64(currentSampleRate - target))
		if diff < minDiff {
			closest = target
			minDiff = diff
		}
	}

	return closest
}

func (g *ExplorerProtoImplV1) getIndices(arr []byte, sep []byte) []int {
	var indices []int
	sepLen := len(sep)
	arrLen := len(arr)

	for i := 0; i <= arrLen-sepLen; i++ {
		if bytes.Equal(arr[i:i+sepLen], sep) {
			indices = append(indices, i)
		}
	}

	return indices
}

func (g *ExplorerProtoImplV1) getChannelData(packetBytes []byte, headerSize, channelSize int) error {
	zOffset := headerSize + int(unsafe.Sizeof(int64(0)))
	zAxisData := make([]int32, channelSize)
	eOffset := zOffset + channelSize*int(unsafe.Sizeof(int32(0)))
	eAxisData := make([]int32, channelSize)
	nOffset := eOffset + channelSize*int(unsafe.Sizeof(int32(0)))
	nAxisData := make([]int32, channelSize)

	recvChecksum := packetBytes[len(packetBytes)-1-3 : len(packetBytes)-1]
	calcChecksum := []uint8{0, 0, 0}
	for i := zOffset; i < eOffset; i++ {
		calcChecksum[0] ^= packetBytes[i]
	}
	for i := eOffset; i < nOffset; i++ {
		calcChecksum[1] ^= packetBytes[i]
	}
	for i := nOffset; i < len(packetBytes)-1-3; i++ {
		calcChecksum[2] ^= packetBytes[i]
	}
	for i := 0; i < len(calcChecksum); i++ {
		if calcChecksum[i] != recvChecksum[i] {
			return fmt.Errorf("checksum mismatch, expected %v, got %v", recvChecksum, calcChecksum)
		}
	}

	err := binary.Read(bytes.NewReader(packetBytes[zOffset:eOffset]), binary.LittleEndian, &zAxisData)
	if err != nil {
		return fmt.Errorf("failed to read z-axis data: %w", err)
	}
	err = binary.Read(bytes.NewReader(packetBytes[eOffset:nOffset]), binary.LittleEndian, &eAxisData)
	if err != nil {
		return fmt.Errorf("failed to read e-axis data: %w", err)
	}
	err = binary.Read(bytes.NewReader(packetBytes[nOffset:len(packetBytes)-1-3]), binary.LittleEndian, &nAxisData)
	if err != nil {
		return fmt.Errorf("failed to read n-axis data: %w", err)
	}

	if len(g.channelDataBuf) != 3 {
		g.channelDataBuf = make([]ChannelData, 3)
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

func (g *ExplorerProtoImplV1) Open(ctx context.Context) (context.Context, context.CancelFunc, error) {
	if g.Transport == nil {
		return nil, nil, errors.New("transport is not opened")
	}
	if g.Logger == nil {
		return nil, nil, errors.New("logger is not set")
	}

	if err := g.Transport.Open(); err != nil {
		return nil, nil, fmt.Errorf("failed to open transport: %w", err)
	}
	ntpClient, err := ntpclient.New(g.NtpOptions.Pool, g.NtpOptions.Retry, g.NtpOptions.ReadTimeout)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ntp client: %w", err)
	}

	g.Logger.Infoln("synchronizing time with NTP servers, it may take a while")
	offset, err := ntpClient.QueryAverage(NTP_MEASUREMENT_ATTEMPTS)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to acquire time from NTP server: %w", err)
	}

	currentTime := time.Now()
	g.TimeSource.Update(currentTime, currentTime.Add(offset))
	g.Logger.Infof("time synchronized with NTP server, local time offset: %d ms", offset.Milliseconds())
	if err = g.Flush(); err != nil {
		return nil, nil, fmt.Errorf("failed to flush transport: %w", err)
	}

	subCtx, cancelFn := context.WithCancel(ctx)

	// In v1 mode, each packet contains 3 channels, n samples per channel.
	// The packet is sent at an interval of (1000 / sample rate) milliseconds.
	// Set n = 5 (also in Explorer) fits the common sample rates (25, 50, 100, 125 Hz).
	const DATA_PACKET_CHANNEL_SIZE = 5

	DATA_PACKET_HEADER := []byte{0xFC, 0x1B}
	packetSize := g.getPacketSize(len(DATA_PACKET_HEADER), DATA_PACKET_CHANNEL_SIZE)
	g.fifoBuffer = fifo.New[byte](10 * packetSize)
	g.messageBus = message.NewBus[EventHandler](EXPLORER_STREAM_TOPIC, 1024)
	g.deviceConfig.SetGnssAvailability(false)

	dummyDeviceId := uint32(0x12F81AC)
	g.deviceVariable.SetDeviceId(&dummyDeviceId)
	g.deviceVariable.SetLatitude(&g.ExplorerOptions.Latitude)
	g.deviceVariable.SetLongitude(&g.ExplorerOptions.Longitude)
	g.deviceVariable.SetElevation(&g.ExplorerOptions.Elevation)
	g.deviceStatus.SetStartedAt(g.TimeSource.Now())
	g.deviceStatus.SetUpdatedAt(time.Unix(0, 0))
	g.deviceConfig.SetProtocol(g.ExplorerOptions.Protocol)
	g.deviceConfig.SetModel(g.ExplorerOptions.Model)

	go func() {
		recvBuf := make([]byte, packetSize)
		prevHeaderIndex := -1

		timeBytes := make([]byte, 8)
		packetBuf := make([]byte, packetSize+len(timeBytes))

		for {
			select {
			case <-subCtx.Done():
				g.Logger.Info("exiting from data packet reader")
				return
			default:
			}

			recvStartTime := g.TimeSource.Now()
			n, err := g.Transport.Read(recvBuf)
			recvEndTime := g.TimeSource.Now()
			if err != nil {
				g.Logger.Errorf("failed to read data from transport: %v", err)
				cancelFn()
			}
			totalLatency := recvEndTime.Sub(recvStartTime) + g.Transport.GetLatency(len(recvBuf))

			// Record the current time of the packet
			currentTime := g.TimeSource.Now().UnixMilli() - totalLatency.Milliseconds()
			binary.BigEndian.PutUint64(timeBytes, uint64(currentTime))

			// Find possible header in the buffer to insert current time next to the header
			headerIndices := g.getIndices(recvBuf[:n], DATA_PACKET_HEADER)
			if len(headerIndices) == 0 {
				continue
			}
			headerIndex := headerIndices[0]
			if prevHeaderIndex == -1 {
				prevHeaderIndex = headerIndex
			}

			// To avoid packet loss, we need to find the "real" header
			// Which is the header that is always equal to the previous header
			for _, index := range headerIndices {
				if index == prevHeaderIndex {
					headerIndex = index
					break
				}
			}
			prevHeaderIndex = headerIndex

			// Copy packet buffer with timestamp
			copy(packetBuf, recvBuf[:headerIndex+len(DATA_PACKET_HEADER)])                                                      // Copy header
			copy(packetBuf[headerIndex+len(DATA_PACKET_HEADER):headerIndex+len(DATA_PACKET_HEADER)+len(timeBytes)], timeBytes)  // Copy timestamp
			copy(packetBuf[headerIndex+len(DATA_PACKET_HEADER)+len(timeBytes):], recvBuf[headerIndex+len(DATA_PACKET_HEADER):]) // Copy packet

			_, _ = g.fifoBuffer.Write(packetBuf...)
		}
	}()

	go func(decodeInterval time.Duration) {
		var (
			expectedNextTimestamp int64
			collectedTimestampArr []int64
		)
		for timer := time.NewTimer(decodeInterval); ; {
			timer.Reset(decodeInterval)

			select {
			case <-timer.C:
				dataPacket, err := g.fifoBuffer.Peek(DATA_PACKET_HEADER, packetSize+8)
				if err != nil {
					continue
				}

				timestamp := int64(binary.BigEndian.Uint64(dataPacket[2:10]))
				if expectedNextTimestamp == 0 {
					expectedNextTimestamp = timestamp
				} else {
					err = g.getChannelData(dataPacket, len(DATA_PACKET_HEADER), DATA_PACKET_CHANNEL_SIZE)
					if err != nil {
						g.Logger.Errorf("failed to get channel data: %v", err)
						g.deviceStatus.IncrementErrors()
						continue
					}
					collectedTimestampArr = append(collectedTimestampArr, timestamp)
				}

				// Calculate proper sample rate to avoid jitter
				currentSampleRate := len(collectedTimestampArr) * DATA_PACKET_CHANNEL_SIZE
				targetSampleRate := g.fixSampleRate(currentSampleRate)

				if math.Abs(float64(timestamp-expectedNextTimestamp)) <= ALLOWED_JITTER_MS_NTP && currentSampleRate == targetSampleRate {
					// Update the next tick even if the buffer is empty
					expectedNextTimestamp = timestamp + time.Second.Milliseconds()
					if len(collectedTimestampArr) == 0 {
						continue
					}

					g.deviceConfig.SetSampleRate(targetSampleRate)
					g.deviceConfig.SetPacketInterval(time.Duration((1000/targetSampleRate)*DATA_PACKET_CHANNEL_SIZE) * time.Millisecond)
					packetTimestamp := collectedTimestampArr[0]
					g.messageBus.Publish(time.UnixMilli(packetTimestamp), &g.deviceConfig, &g.deviceVariable, g.channelDataBuf)
					g.deviceStatus.IncrementMessages()

					collectedTimestampArr = []int64{}
					g.channelDataBuf = []ChannelData{}
				} else if timestamp-expectedNextTimestamp > ALLOWED_JITTER_MS_NTP {
					g.Logger.Warnf("jitter detected, discarding this packet, expected %v, got %v", expectedNextTimestamp, timestamp)
					g.deviceStatus.IncrementErrors()
					// Update the next tick, clear the buffer if the jitter exceeds the threshold
					expectedNextTimestamp = timestamp + time.Second.Milliseconds()
					collectedTimestampArr = []int64{}
					g.channelDataBuf = []ChannelData{}
				}

				g.deviceStatus.IncrementFrames()
				g.deviceStatus.SetUpdatedAt(time.UnixMilli(timestamp))
			case <-subCtx.Done():
				g.Logger.Info("exiting from data packet decoder")
				timer.Stop()
				return
			}
		}
	}(5 * time.Millisecond)

	go func(resyncInterval time.Duration) {
		for timer := time.NewTimer(resyncInterval); ; {
			select {
			case <-timer.C:
				g.Logger.Info("re-synchronizing time with NTP servers")
				offset, err := ntpClient.QueryAverage(NTP_MEASUREMENT_ATTEMPTS)
				if err != nil {
					g.Logger.Warnf("error occurred while re-synchronizing time with NTP: %v", err)
					timer.Reset(resyncInterval)
					continue
				}
				timer.Reset(resyncInterval)
				currentTime := time.Now()
				g.TimeSource.Update(currentTime, currentTime.Add(offset))
				g.Logger.Infof("time synchronized with NTP server, local time offset: %d ms", offset.Milliseconds())
			case <-subCtx.Done():
				timer.Stop()
				return
			}
		}
	}(NTP_RESYNC_INTERVAL)

	return subCtx, cancelFn, nil
}

func (g *ExplorerProtoImplV1) Close() error {
	if g.Transport == nil {
		return errors.New("transport is not opened")
	}

	return g.Transport.Close()
}

func (g *ExplorerProtoImplV1) Subscribe(clientId string, handler EventHandler) error {
	return g.messageBus.Subscribe(clientId, handler)
}

func (g *ExplorerProtoImplV1) Unsubscribe(clientId string) error {
	return g.messageBus.Unsubscribe(clientId)
}

func (g *ExplorerProtoImplV1) GetConfig() DeviceConfig {
	return DeviceConfig{
		packetInterval: g.deviceConfig.GetPacketInterval(),
		sampleRate:     g.deviceConfig.GetSampleRate(),
		gnssEnabled:    g.deviceConfig.GetGnssAvailability(),
		channelCodes:   g.deviceConfig.GetChannelCodes(),
		model:          g.deviceConfig.GetModel(),
		protocol:       g.deviceConfig.GetProtocol(),
	}
}

func (g *ExplorerProtoImplV1) GetStatus() DeviceStatus {
	return DeviceStatus{
		startedAt: g.deviceStatus.GetStartedAt(),
		updatedAt: g.deviceStatus.GetUpdatedAt(),
		frames:    g.deviceStatus.GetFrames(),
		errors:    g.deviceStatus.GetErrors(),
		messages:  g.deviceStatus.GetMessages(),
	}
}

func (g *ExplorerProtoImplV1) GetCoordinates(fuzzy bool) (float64, float64, float64, error) {
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

func (g *ExplorerProtoImplV1) GetTemperature() (float64, error) {
	temp, err := g.deviceVariable.GetTemperature()
	if err != nil {
		return 0, fmt.Errorf("failed to get temperature: %w", err)
	}
	return temp, nil
}

func (g *ExplorerProtoImplV1) GetDeviceId() string {
	devId, err := g.deviceVariable.GetDeviceId()
	if err != nil {
		return "N/A"
	}
	return fmt.Sprintf("%08X", devId)
}

func (g *ExplorerProtoImplV1) Flush() error {
	return g.Transport.Flush()
}

func (g *ExplorerProtoImplV1) GetMetadata(stationAffiliation, stationDescription, stationCountry, stationPlace, networkCode, stationCode, locationCode string, fuzzyCoordinates bool) (metadata.IMetadata, error) {
	latitude, longitude, elevation, err := g.GetCoordinates(fuzzyCoordinates)
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
