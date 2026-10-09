package explorer

import (
	"context"
	"encoding/binary"
	"time"

	"github.com/anyshake/observer/pkg/timesource"
)

func (g *ExplorerProtoImplV1) readStream(subCtx context.Context, packetSize int) {
	defer close(g.streamDone)
	defer g.cancelFn()
	const DATA_PACKET_CHANNEL_SIZE = 5
	DATA_PACKET_HEADER := []byte{0xFC, 0x1B}
	recvBuf := make([]byte, packetSize)
	prevHeaderIndex := -1

	timeBytes := make([]byte, 8)
	packetBuf := make([]byte, packetSize+len(timeBytes))

	for {
		select {
		case <-subCtx.Done():
			g.Logger.Infoln("exiting from data packet reader")
			return
		default:
		}

		recvStartTime := g.TimeSource.Now()
		n, err := g.Transport.Read(recvBuf)
		recvEndTime := g.TimeSource.Now()
		if err != nil {
			g.Logger.Errorf("failed to read data from transport: %v", err)
			return
		}
		if subCtx.Err() != nil {
			return
		}
		elapsed := recvEndTime.Sub(recvStartTime)
		latency := g.Transport.GetLatency(len(recvBuf))

		// Calculate proper sample rate to avoid jitter
		currentSampleRate, err := g.fixSampleRate(DATA_PACKET_CHANNEL_SIZE, elapsed)
		if err != nil {
			g.Logger.Errorf("failed to determine current sample rate: %v", err)
			continue
		}
		g.deviceConfig.SetSampleRate(currentSampleRate)
		g.deviceConfig.SetPacketInterval(time.Duration(1000/currentSampleRate*DATA_PACKET_CHANNEL_SIZE) * time.Millisecond)

		// Record the current time of the packet
		currentTime := g.TimeSource.Now().UnixMilli() - (elapsed + latency).Milliseconds()
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
}

func (g *ExplorerProtoImplV1) decodeStream(subCtx context.Context, packetSize int) {
	defer close(g.decoderDone)
	const DATA_PACKET_CHANNEL_SIZE = 5
	DATA_PACKET_HEADER := []byte{0xFC, 0x1B}
	const decodeInterval = 5 * time.Millisecond
	var (
		collectedTimestampArr []int64
	)
	for timer := time.NewTimer(decodeInterval); ; {
		timer.Reset(decodeInterval)

		select {
		case <-timer.C:
			dataPacket, err := g.fifoBuffer.Peek(DATA_PACKET_HEADER, packetSize+8) // extra 8 bytes for inserting timestamp
			if err != nil {
				continue
			}

			currentSampleRate := g.deviceConfig.GetSampleRate()
			if currentSampleRate > 0 {
				timestamp := int64(binary.BigEndian.Uint64(dataPacket[2:10]))
				channelData, err := g.getChannelData(dataPacket, len(DATA_PACKET_HEADER), DATA_PACKET_CHANNEL_SIZE)
				if err != nil {
					g.Logger.Errorf("failed to get channel data: %v", err)
					g.deviceStatus.IncrementErrors()
					continue
				}

				collectedTimestampArr = append(collectedTimestampArr, timestamp)
				g.deviceStatus.IncrementFrames()

				g.messageBusRealtime.Publish(NewEvent(time.UnixMilli(timestamp), &g.deviceConfig, channelData))
				if len(collectedTimestampArr)*DATA_PACKET_CHANNEL_SIZE == currentSampleRate {
					packetTimestamp := collectedTimestampArr[0]
					g.messageBus.Publish(NewEvent(time.UnixMilli(packetTimestamp), &g.deviceConfig, g.channelDataBuf))
					g.deviceStatus.IncrementMessages()
					collectedTimestampArr = []int64{}
					g.channelDataBuf = []ChannelData{}
				} else if len(collectedTimestampArr)*DATA_PACKET_CHANNEL_SIZE > currentSampleRate {
					g.Logger.Warnf("packet timestamp is not in sync with current sample rate, packet timestamp: %v, current sample rate: %v", collectedTimestampArr[0], currentSampleRate)
					collectedTimestampArr = []int64{}
					g.channelDataBuf = []ChannelData{}
					g.deviceStatus.IncrementErrors()
				}

				g.deviceStatus.SetUpdatedAt(time.UnixMilli(timestamp))
			}
		case <-subCtx.Done():
			g.Logger.Infoln("exiting from data packet decoder")
			timer.Stop()
			return
		}
	}
}

func (g *ExplorerProtoImplV1) synchronizeNTP(subCtx context.Context) {
	const resyncInterval = NTP_RESYNC_INTERVAL
	defer close(g.ntpDone)
	timer := time.NewTimer(resyncInterval)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			g.Logger.Infoln("re-synchronizing time with NTP servers")
			offset, server, err := g.ntpClient.QueryContext(subCtx)
			if subCtx.Err() != nil {
				return
			}
			if err != nil {
				g.Logger.Warnf("error occurred while re-synchronizing time with NTP: %v", err)
				timer.Reset(resyncInterval)
				continue
			}
			timer.Reset(resyncInterval)
			currentMonotonicTime := timesource.MonotonicNow()
			g.clockDriftBuf.Push(clockDrift{offset: offset, measuredAt: currentMonotonicTime})
			ppm := getLongTermClockDriftPPM(g.clockDriftBuf, NTP_PPM_MEASURE_WINDOW)
			g.TimeSource.Update(currentMonotonicTime, currentMonotonicTime.Add(offset), ppm, nil)
			g.Logger.Infof("time synchronized with NTP server: %s, local monotonic time offset: %d ms, clock drift PPM: %.2f", server, offset.Milliseconds(), ppm)
		case <-subCtx.Done():
			return
		}
	}
}
