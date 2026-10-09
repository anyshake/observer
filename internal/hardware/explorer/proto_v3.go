package explorer

import (
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
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/message"
	"github.com/anyshake/observer/pkg/metadata"
	"github.com/anyshake/observer/pkg/ntpclient"
	"github.com/anyshake/observer/pkg/ringbuf"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/anyshake/observer/pkg/transport"
)

type ExplorerProtoImplV3 struct {
	ChannelCodes    []string
	ExplorerOptions ExplorerOptions
	NtpOptions      NtpOptions
	Logger          *logger.Adapter
	TimeSource      *timesource.Source

	Transport   transport.ITransport
	fifoBuffer  *fifo.Buffer[*explorerProtocolPacketV3]
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

	prevMcuTimestamp    int64
	isDataStreamStable  bool
	timeCalibrationChan chan [2]time.Time

	// Protects reader/decoder state transitions and archive buffers.
	streamMutex sync.Mutex
	stream      explorerStreamV3

	flagMutex        sync.Mutex
	variableAllSet   bool
	collectedSamples int

	packetTimeObj  time.Time
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

func (g *ExplorerProtoImplV3) variablesReady() bool {
	g.flagMutex.Lock()
	defer g.flagMutex.Unlock()
	return g.variableAllSet
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

	g.flagMutex.Lock()
	g.variableAllSet = variableAllSet
	g.flagMutex.Unlock()
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
				v := int32(channelDataBytes[offset]) | int32(channelDataBytes[offset+1])<<8 | int32(channelDataBytes[offset+2])<<16
				ch.Data[i] = (v << 8) >> 8 // sign extend for 24-bit
				offset += 3
			case "int32":
				ch.Data[i] = int32(binary.LittleEndian.Uint32(channelDataBytes[offset : offset+4]))
				offset += 4
			}
		}
	}
}

func (g *ExplorerProtoImplV3) verifyChecksum(packetData []byte) error {
	if len(packetData) == 0 {
		return errors.New("empty packet data")
	}

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
	opened := false
	defer func() {
		if !opened {
			_ = g.Close()
		}
	}()
	if err := g.Flush(); err != nil {
		return nil, nil, fmt.Errorf("failed to flush transport: %w", err)
	}

	if g.Logger == nil {
		return nil, nil, errors.New("logger is not set")
	}
	ntpClient, err := ntpclient.New(g.NtpOptions.Pool, g.NtpOptions.Retry, g.NtpOptions.ReadTimeout, timesource.MonotonicNow, ntpclient.WithLogger(g.Logger))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ntp client: %w", err)
	}

	subCtx, cancelFn := context.WithCancel(ctx)
	g.ntpClient, g.cancelFn = ntpClient, cancelFn

	g.fifoBuffer = fifo.New[*explorerProtocolPacketV3](512)
	g.clockDriftBuf = ringbuf.New[clockDrift](100)
	g.messageBus = message.NewBus[Event](EXPLORER_STREAM_TOPIC)
	g.messageBusRealtime = message.NewBus[Event](EXPLORER_REALTIME_STREAM_TOPIC)
	g.deviceStatus.SetUpdatedAt(time.Unix(0, 0))
	g.deviceConfig.SetProtocol(g.ExplorerOptions.Protocol)
	g.deviceConfig.SetModel(filepath.Base(g.ExplorerOptions.Model))
	g.timeCalibrationChan = make(chan [2]time.Time)

	var initFlag int32
	atomic.StoreInt32(&initFlag, 0)
	readyChan := make(chan struct{})

	g.streamDone = make(chan struct{})
	g.decoderDone = make(chan struct{})
	g.ntpDone = make(chan struct{})
	g.stream = explorerStreamV3{}
	g.prevMcuTimestamp = 0
	g.isDataStreamStable = false
	g.resetVariables()
	g.resetFlags()
	go g.readStream(subCtx, readyChan, &initFlag)
	go g.decodeStream(subCtx, &initFlag)
	go g.synchronizeNTP(subCtx, readyChan)

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

func (g *ExplorerProtoImplV3) Close() error {
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

func (g *ExplorerProtoImplV3) Subscribe(clientId string, handler EventHandler) error {
	return g.messageBus.Subscribe(clientId, normalStreamSubscriptionOptions(func(err error) {
		if !errors.Is(err, message.ErrBusClosed) {
			g.Logger.Errorf("normal stream subscriber %s failed: %v", clientId, err)
		}
	}), handler)
}

func (g *ExplorerProtoImplV3) Unsubscribe(clientId string) error {
	return g.messageBus.Unsubscribe(clientId)
}

func (g *ExplorerProtoImplV3) SubscribeRealtime(clientId string, handler EventHandler) error {
	return g.messageBusRealtime.Subscribe(clientId, realtimeStreamSubscriptionOptions(func(err error) {
		if !errors.Is(err, message.ErrBusClosed) {
			g.Logger.Errorf("realtime stream subscriber %s failed: %v", clientId, err)
		}
	}), handler)
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
