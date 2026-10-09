package explorer

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"time"
	"unsafe"

	"github.com/anyshake/observer/pkg/fifo"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/message"
	"github.com/anyshake/observer/pkg/metadata"
	"github.com/anyshake/observer/pkg/ntpclient"
	"github.com/anyshake/observer/pkg/ringbuf"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/anyshake/observer/pkg/transport"
)

type ExplorerProtoImplV1 struct {
	ChannelCodes    []string
	ExplorerOptions ExplorerOptions
	NtpOptions      NtpOptions
	Logger          *logger.Adapter
	TimeSource      *timesource.Source

	Transport   transport.ITransport
	fifoBuffer  *fifo.Buffer[byte]
	ntpClient   *ntpclient.Client
	cancelFn    context.CancelFunc
	ntpDone     chan struct{}
	streamDone  chan struct{}
	decoderDone chan struct{}

	// buf length: 100; ppm window: 60 min
	clockDriftBuf *ringbuf.Buffer[clockDrift]

	// 1 message per second, for archiving service, etc.
	messageBus *message.Bus[Event]
	// 1 message per packet, for realtime purposes
	messageBusRealtime *message.Bus[Event]

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

func (g *ExplorerProtoImplV1) fixSampleRate(channelSize int64, duration time.Duration) (int, error) {
	if duration.Milliseconds() == 0 {
		return 0, errors.New("invalid duration")
	}

	currentSampleRate := int(1000 / duration.Milliseconds() * channelSize)
	currentSampleRate = int(math.Round(float64(currentSampleRate)/5.0) * 5.0)

	// All divisors of 5000 greater than or equal to 5
	targetSampleRates := []int{50, 100, 125, 200, 250, 500, 1000, 1250, 2500, 5000}
	closest := targetSampleRates[0]
	minDiff := math.Abs(float64(currentSampleRate - closest))

	for _, target := range targetSampleRates {
		diff := math.Abs(float64(currentSampleRate - target))
		if diff < minDiff {
			closest = target
			minDiff = diff
		}
	}

	return closest, nil
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

func (g *ExplorerProtoImplV1) getChannelData(packetBytes []byte, headerSize, channelSize int) (channelData []ChannelData, err error) {
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
			return nil, fmt.Errorf("checksum mismatch, expected %v, got %v", recvChecksum, calcChecksum)
		}
	}

	if err = binary.Read(bytes.NewReader(packetBytes[zOffset:eOffset]), binary.LittleEndian, &zAxisData); err != nil {
		return nil, fmt.Errorf("failed to read z-axis data: %w", err)
	}
	if err = binary.Read(bytes.NewReader(packetBytes[eOffset:nOffset]), binary.LittleEndian, &eAxisData); err != nil {
		return nil, fmt.Errorf("failed to read e-axis data: %w", err)
	}
	if err = binary.Read(bytes.NewReader(packetBytes[nOffset:len(packetBytes)-1-3]), binary.LittleEndian, &nAxisData); err != nil {
		return nil, fmt.Errorf("failed to read n-axis data: %w", err)
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
	opened := false
	defer func() {
		if !opened {
			_ = g.Close()
		}
	}()
	ntpClient, err := ntpclient.New(g.NtpOptions.Pool, g.NtpOptions.Retry, g.NtpOptions.ReadTimeout, timesource.MonotonicNow, ntpclient.WithLogger(g.Logger))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ntp client: %w", err)
	}
	subCtx, cancelFn := context.WithCancel(ctx)
	g.ntpClient, g.cancelFn = ntpClient, cancelFn

	g.Logger.Infoln("synchronizing time with NTP servers, it may take a while")
	offset, err := ntpClient.QueryAverageContext(subCtx, NTP_MEASUREMENT_ATTEMPTS)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to acquire time from NTP server: %w", err)
	}

	currentMonotonicTime := timesource.MonotonicNow()
	g.TimeSource.Update(currentMonotonicTime, currentMonotonicTime.Add(offset), 0, timesource.MonotonicNow)
	g.Logger.Infof("time synchronized with NTP server, local monotonic time offset: %d ms", offset.Milliseconds())
	if err = g.Flush(); err != nil {
		return nil, nil, fmt.Errorf("failed to flush transport: %w", err)
	}

	// In v1 mode, each packet contains 3 channels, n samples per channel.
	// The packet is sent at an interval of (1000 / sample rate) milliseconds.
	// Set n = 5 (also in Explorer) fits the common sample rates (25, 50, 100, 125 Hz).
	const DATA_PACKET_CHANNEL_SIZE = 5

	DATA_PACKET_HEADER := []byte{0xFC, 0x1B}
	packetSize := g.getPacketSize(len(DATA_PACKET_HEADER), DATA_PACKET_CHANNEL_SIZE)
	g.fifoBuffer = fifo.New[byte](10 * packetSize)
	g.clockDriftBuf = ringbuf.New[clockDrift](100)
	g.messageBus = message.NewBus[Event](EXPLORER_STREAM_TOPIC)
	g.messageBusRealtime = message.NewBus[Event](EXPLORER_REALTIME_STREAM_TOPIC)
	g.deviceConfig.SetGnssAvailability(false)

	dummyDeviceId := uint32(0x12F81AC)
	g.deviceVariable.SetDeviceId(&dummyDeviceId)
	g.deviceVariable.SetLatitude(&g.ExplorerOptions.Latitude)
	g.deviceVariable.SetLongitude(&g.ExplorerOptions.Longitude)
	g.deviceVariable.SetElevation(&g.ExplorerOptions.Elevation)
	g.deviceStatus.SetStartedAt(g.TimeSource.Now())
	g.deviceStatus.SetUpdatedAt(time.Unix(0, 0))
	g.deviceConfig.SetProtocol(g.ExplorerOptions.Protocol)
	g.deviceConfig.SetModel(filepath.Base(g.ExplorerOptions.Model))

	g.streamDone = make(chan struct{})
	g.decoderDone = make(chan struct{})
	g.ntpDone = make(chan struct{})
	go g.readStream(subCtx, packetSize)
	go g.decodeStream(subCtx, packetSize)
	go g.synchronizeNTP(subCtx)

	opened = true
	return subCtx, cancelFn, nil
}

func (g *ExplorerProtoImplV1) Close() error {
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
	if g.decoderDone != nil {
		<-g.decoderDone
	}
	if g.messageBus != nil {
		g.messageBus.Close()
	}
	if g.messageBusRealtime != nil {
		g.messageBusRealtime.Close()
	}

	return err
}

func (g *ExplorerProtoImplV1) Subscribe(clientId string, handler EventHandler) error {
	return g.messageBus.Subscribe(clientId, normalStreamSubscriptionOptions(func(err error) {
		if !errors.Is(err, message.ErrBusClosed) {
			g.Logger.Errorf("normal stream subscriber %s failed: %v", clientId, err)
		}
	}), handler)
}

func (g *ExplorerProtoImplV1) Unsubscribe(clientId string) error {
	return g.messageBus.Unsubscribe(clientId)
}

func (g *ExplorerProtoImplV1) SubscribeRealtime(clientId string, handler EventHandler) error {
	return g.messageBusRealtime.Subscribe(clientId, realtimeStreamSubscriptionOptions(func(err error) {
		if !errors.Is(err, message.ErrBusClosed) {
			g.Logger.Errorf("realtime stream subscriber %s failed: %v", clientId, err)
		}
	}), handler)
}

func (g *ExplorerProtoImplV1) UnsubscribeRealtime(clientId string) error {
	return g.messageBusRealtime.Unsubscribe(clientId)
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

func (g *ExplorerProtoImplV1) GetMetadata(stationAffiliation, stationDescription, stationCountry, stationPlace, networkCode, stationCode, locationCode string, fuzzyCoordinates bool) (*metadata.Render, error) {
	latitude, longitude, elevation, err := g.GetCoordinates(fuzzyCoordinates)
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
