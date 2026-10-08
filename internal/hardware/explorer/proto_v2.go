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
	"time"
	"unsafe"

	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/message"
	"github.com/anyshake/observer/pkg/metadata"
	"github.com/anyshake/observer/pkg/ntpclient"
	"github.com/anyshake/observer/pkg/ringbuf"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/anyshake/observer/pkg/transport"
)

type ExplorerProtoImplV2 struct {
	ChannelCodes    []string
	ExplorerOptions ExplorerOptions
	NtpOptions      NtpOptions
	Logger          *logger.Adapter
	TimeSource      *timesource.Source

	Transport  transport.ITransport
	ntpClient  *ntpclient.Client
	cancelFn   context.CancelFunc
	ntpDone    chan struct{}
	streamDone chan struct{}

	// buf length: 100; ppm window: 60 min
	clockDriftBuf *ringbuf.Buffer[clockDrift]

	// 1 message per second, for archiving service, etc.
	messageBus *message.Bus[Event]
	// 1 message per packet, for realtime purposes
	messageBusRealtime *message.Bus[Event]

	// Serializes packet processing and clock updates from the NTP worker.
	streamMutex    sync.Mutex
	stream         explorerStreamV2
	variableAllSet bool

	deviceStatus   DeviceStatus
	deviceConfig   DeviceConfig
	deviceVariable DeviceVariable
	channelDataBuf []ChannelData
}

func (g *ExplorerProtoImplV2) resetVariables() {
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

func (g *ExplorerProtoImplV2) variablesReady() bool {
	g.streamMutex.Lock()
	defer g.streamMutex.Unlock()
	return g.variableAllSet
}

func (g *ExplorerProtoImplV2) getVariableData(mcuTimestamp int64, variableBytes uint32) {
	gnssEnabled := g.deviceConfig.GetGnssAvailability()
	switch (mcuTimestamp / 1000) % 4 {
	case 0:
		gnssEnable := (variableBytes&0x80000000 != 0)
		if gnssEnable != gnssEnabled {
			g.resetVariables()
		}
		gnssEnabled = gnssEnable
		deviceId := variableBytes & 0x7FFFFFFF
		g.deviceVariable.SetDeviceId(&deviceId)
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

func (g *ExplorerProtoImplV2) getChannelData(packetBytes []byte, headerSize, channelSize int) (channelData []ChannelData, err error) {
	if len(g.channelDataBuf) != 3 {
		g.channelDataBuf = make([]ChannelData, 3)
	}

	zOffset := headerSize + int(unsafe.Sizeof(int64(0))) + int(unsafe.Sizeof(uint32(0)))
	zAxisData := make([]int32, channelSize)
	eOffset := zOffset + (channelSize)*int(unsafe.Sizeof(int32(0)))
	eAxisData := make([]int32, channelSize)
	nOffset := eOffset + (channelSize)*int(unsafe.Sizeof(int32(0)))
	nAxisData := make([]int32, channelSize)

	if err := binary.Read(bytes.NewReader(packetBytes[zOffset:eOffset]), binary.LittleEndian, &zAxisData); err != nil {
		return nil, fmt.Errorf("failed to read z-axis data: %w", err)
	}
	if err = binary.Read(bytes.NewReader(packetBytes[eOffset:nOffset]), binary.LittleEndian, &eAxisData); err != nil {
		return nil, fmt.Errorf("failed to read e-axis data: %w", err)
	}
	if err = binary.Read(bytes.NewReader(packetBytes[nOffset:len(packetBytes)-1]), binary.LittleEndian, &nAxisData); err != nil {
		return nil, fmt.Errorf("failed to read n-axis data: %w", err)
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

	result := make([]ChannelData, 3)
	result[0] = ChannelData{
		ChannelCode: g.channelDataBuf[0].ChannelCode,
		ChannelId:   1,
		ByteSize:    4,
		DataType:    "int32",
		Data:        zAxisData,
	}
	result[1] = ChannelData{
		ChannelCode: g.channelDataBuf[1].ChannelCode,
		ChannelId:   2,
		ByteSize:    4,
		DataType:    "int32",
		Data:        eAxisData,
	}
	result[2] = ChannelData{
		ChannelCode: g.channelDataBuf[2].ChannelCode,
		ChannelId:   3,
		ByteSize:    4,
		DataType:    "int32",
		Data:        nAxisData,
	}
	return result, nil
}

func (g *ExplorerProtoImplV2) verifyChecksum(packetData, header []byte) error {
	if len(packetData) == 0 {
		return errors.New("empty packet data")
	}

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
	opened := false
	defer func() {
		if !opened {
			_ = g.Close()
		}
	}()
	if err := g.Flush(); err != nil {
		return nil, nil, fmt.Errorf("failed to flush transport: %w", err)
	}
	ntpClient, err := ntpclient.New(g.NtpOptions.Pool, g.NtpOptions.Retry, g.NtpOptions.ReadTimeout, timesource.MonotonicNow, ntpclient.WithLogger(g.Logger))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ntp client: %w", err)
	}

	subCtx, cancelFn := context.WithCancel(ctx)
	g.ntpClient, g.cancelFn = ntpClient, cancelFn

	g.clockDriftBuf = ringbuf.New[clockDrift](100)
	g.messageBus = message.NewBus[Event](EXPLORER_STREAM_TOPIC)
	g.messageBusRealtime = message.NewBus[Event](EXPLORER_REALTIME_STREAM_TOPIC)
	g.deviceStatus.SetUpdatedAt(time.Unix(0, 0))
	g.deviceConfig.SetProtocol(g.ExplorerOptions.Protocol)
	g.deviceConfig.SetModel(filepath.Base(g.ExplorerOptions.Model))
	g.streamMutex.Lock()
	g.stream = explorerStreamV2{}
	g.resetStream()
	g.streamMutex.Unlock()
	g.streamDone = make(chan struct{})
	g.ntpDone = make(chan struct{})
	readyChan := make(chan struct{})
	go g.readStream(subCtx, readyChan)
	go g.synchronizeNTP(subCtx)

	select {
	case <-readyChan:
	case <-subCtx.Done():
	}
	if err := subCtx.Err(); err != nil {
		return nil, nil, err
	}
	opened = true
	return subCtx, cancelFn, nil
}

func (g *ExplorerProtoImplV2) Close() error {
	if g.cancelFn != nil {
		g.cancelFn()
	}
	if g.ntpClient != nil {
		_ = g.ntpClient.Close()
	}
	if g.ntpDone != nil {
		<-g.ntpDone
	}
	if g.Transport == nil {
		return errors.New("transport is not opened")
	}
	err := g.Transport.Close()
	if g.streamDone != nil {
		<-g.streamDone
	}
	if g.messageBus != nil {
		g.messageBus.Close()
	}
	if g.messageBusRealtime != nil {
		g.messageBusRealtime.Close()
	}

	return err
}

func (g *ExplorerProtoImplV2) Subscribe(clientId string, handler EventHandler) error {
	return g.messageBus.Subscribe(clientId, normalStreamSubscriptionOptions(func(err error) {
		if !errors.Is(err, message.ErrBusClosed) {
			g.Logger.Errorf("normal stream subscriber %s failed: %v", clientId, err)
		}
	}), handler)
}

func (g *ExplorerProtoImplV2) Unsubscribe(clientId string) error {
	return g.messageBus.Unsubscribe(clientId)
}

func (g *ExplorerProtoImplV2) SubscribeRealtime(clientId string, handler EventHandler) error {
	return g.messageBusRealtime.Subscribe(clientId, realtimeStreamSubscriptionOptions(func(err error) {
		if !errors.Is(err, message.ErrBusClosed) {
			g.Logger.Errorf("realtime stream subscriber %s failed: %v", clientId, err)
		}
	}), handler)
}

func (g *ExplorerProtoImplV2) UnsubscribeRealtime(clientId string) error {
	return g.messageBusRealtime.Unsubscribe(clientId)
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

func (g *ExplorerProtoImplV2) GetMetadata(stationAffiliation, stationDescription, stationCountry, stationPlace, networkCode, stationCode, locationCode string, fuzzyLocation bool) (*metadata.Render, error) {
	latitude, longitude, elevation, err := g.GetCoordinates(fuzzyLocation)
	if err != nil {
		return nil, err
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
