package explorer

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/anyshake/observer/internal/hardware/explorer/metadata"
	"github.com/anyshake/observer/pkg/fifo"
	"github.com/anyshake/observer/pkg/message"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/anyshake/observer/pkg/transport"
	"github.com/sirupsen/logrus"
)

type ExplorerProtoImplV3 struct {
	ChannelCodes []string
	Model        string
	Logger       *logrus.Entry
	TimeSource   *timesource.Source

	FallbackLatitude  float64
	FallbackLongitude float64
	FallbackElevation float64

	Transport  transport.ITransport
	fifoBuffer fifo.Buffer[byte]
	messageBus message.Bus[EventHandler]

	prevTimestamp4NonGnssMode uint64
	timeDiff4NonGnssMode      uint64

	packetTimeObj    time.Time
	variableAllSet   bool
	collectedSamples int

	deviceStatus   DeviceStatus
	deviceConfig   DeviceConfig
	deviceVariable DeviceVariable
	channelDataBuf []ChannelData
}

func (g *ExplorerProtoImplV3) resetFlags() {
	g.channelDataBuf = []ChannelData{}
	g.packetTimeObj = time.Time{}
	g.collectedSamples = 0

}

func (g *ExplorerProtoImplV3) resetVariables() {
	g.variableAllSet = false
	g.deviceVariable.Reset()
}

func (g *ExplorerProtoImplV3) getTimestamp(mcuTimestamp uint64) uint64 {
	if g.deviceConfig.GetGnssAvailability() {
		return mcuTimestamp
	}

	currentTime := g.TimeSource.Get()
	currentUnixMilli := currentTime.UnixMilli()

	// Force calibrate time difference every day UTC at 00:00:00
	currentUtcDate := time.Unix(currentUnixMilli/1000, 0).UTC().Format("2006-01-02")
	prevUtcDate := time.Unix(int64(g.prevTimestamp4NonGnssMode)/1000, 0).UTC().Format("2006-01-02")
	if g.timeDiff4NonGnssMode == 0 || currentUtcDate != prevUtcDate {
		g.timeDiff4NonGnssMode = uint64(currentUnixMilli) - mcuTimestamp
	}

	timestamp := mcuTimestamp + g.timeDiff4NonGnssMode

	// In this case the hardware may have been reset, we need to update the time difference
	// Reset the variableAllSet flag to false, clear all the pointers in deviceInfo
	if timestamp < g.prevTimestamp4NonGnssMode {
		g.timeDiff4NonGnssMode = uint64(currentUnixMilli) - mcuTimestamp
		timestamp = mcuTimestamp + g.timeDiff4NonGnssMode
		g.resetVariables()
		g.resetFlags()
	}

	g.prevTimestamp4NonGnssMode = timestamp
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

func (g *ExplorerProtoImplV3) getVariableData(mcuTimestamp uint64, deviceConfig uint32, variableBytes [4]byte) {
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
			g.deviceVariable.SetLatitude(&g.FallbackLatitude)
		}
	case 2:
		if (variableBits>>2)&0x1 != 0 {
			n := float64(math.Float32frombits(variableData))
			g.deviceVariable.SetLongitude(&n)
		} else {
			g.deviceVariable.SetLongitude(&g.FallbackLongitude)
		}
	case 3:
		if (variableBits>>3)&0x1 != 0 {
			n := float64(math.Float32frombits(variableData))
			g.deviceVariable.SetElevation(&n)
		} else {
			g.deviceVariable.SetElevation(&g.FallbackElevation)
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
		if _, err := g.deviceVariable.GetLatitude(); err != nil {
			variableAllSet = false
		}
	}

	if (variableBits>>2)&0x1 != 0 {
		if _, err := g.deviceVariable.GetLongitude(); err != nil {
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

	subCtx, cancelFn := context.WithCancel(ctx)

	g.fifoBuffer = fifo.New[byte](4096)
	g.messageBus = message.NewBus[EventHandler](EXPLORER_STREAM_TOPIC, 1024)
	g.deviceStatus.SetStartedAt(g.TimeSource.Get())
	g.deviceStatus.SetUpdatedAt(time.Unix(0, 0))
	g.deviceConfig.SetProtocol("v3")
	g.deviceConfig.SetModel(g.Model)

	// FIXME: Timestamp drift if software starts before hardware
	go func(readInterval time.Duration) {
		timer := time.NewTimer(readInterval)
		buf := make([]byte, 512)

		for {
			timer.Reset(readInterval)

			select {
			case <-timer.C:
				n, err := g.Transport.Read(buf)
				if err != nil {
					g.Logger.Errorf("failed to read data from transport: %v", err)
					cancelFn()
				}
				_, _ = g.fifoBuffer.Write(buf[:n]...)
			case <-subCtx.Done():
				g.Logger.Info("exiting from data packet reader")
				timer.Stop()
				return
			}
		}
	}(10 * time.Millisecond)

	go func(decodeInterval time.Duration) {
		DATA_PACKET_HEADER := []byte{0x01, 0xFE}
		DATA_PACKET_TAILER := []byte{0xEF, 0x10}

		timer := time.NewTimer(decodeInterval)

		for {
			timer.Reset(decodeInterval)

			select {
			case <-timer.C:
				// Fixed 18 bytes for header + timestamp + device config + variable data
				packetFixedSection, err := g.fifoBuffer.Peek(DATA_PACKET_HEADER, len(DATA_PACKET_HEADER)+16)
				if err != nil {
					continue
				}

				mcuTimestamp := binary.LittleEndian.Uint64(packetFixedSection[2:10])
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
				for !g.fifoBuffer.Available(readSize) {
					time.Sleep(time.Millisecond)
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

				calcChecksum := uint8(0)
				for _, b := range packetFixedSection[len(DATA_PACKET_HEADER):] {
					calcChecksum ^= b
				}
				for _, b := range channelDataSection[:len(channelDataSection)-1-len(recvTailer)] {
					calcChecksum ^= b
				}
				recvChecksum := channelDataSection[len(channelDataSection)-len(DATA_PACKET_TAILER)-1]
				if recvChecksum != calcChecksum {
					g.Logger.Errorf("checksum mismatch, expected %v, got %v", recvChecksum, calcChecksum)
					g.deviceStatus.IncrementErrors()
					continue
				}

				g.deviceStatus.IncrementFrames()
				g.deviceStatus.SetUpdatedAt(timeObj)

				if !g.variableAllSet {
					g.Logger.Infoln("waiting for device config to be fully collected, this may take a while")
					continue
				}

				g.getChannelData(channelData, channelDataSection[:len(channelDataSection)-1-len(recvTailer)], channelChunkLength)
				g.collectedSamples += channelChunkLength

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
				if g.collectedSamples < sampleRate {
					continue
				} else if g.collectedSamples == sampleRate {
					g.messageBus.Publish(timeObj, &g.deviceConfig, &g.deviceVariable, g.channelDataBuf)
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
	}(50 * time.Millisecond)

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

func (g *ExplorerProtoImplV3) GetCoordinates() (float64, float64, float64, error) {
	lat, err := g.deviceVariable.GetLatitude()
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get latitude: %w", err)
	}

	lon, err := g.deviceVariable.GetLongitude()
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

func (g *ExplorerProtoImplV3) GetMetadata(stationAffiliation, stationDescription, stationCountry, stationPlace, networkCode, stationCode, locationCode string) (metadata.IMetadata, error) {
	latitude, err := g.deviceVariable.GetLatitude()
	if err != nil {
		return nil, fmt.Errorf("failed to get latitude: %w", err)
	}
	longitude, err := g.deviceVariable.GetLongitude()
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
