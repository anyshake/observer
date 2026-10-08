package explorer

import (
	"bytes"
	"context"
	"encoding/binary"
	"math"
	"sync/atomic"
	"time"

	"github.com/anyshake/observer/pkg/timesource"
	"github.com/samber/lo"
)

type explorerProtocolPacketV3 struct {
	dataTime  time.Time
	dataBytes []byte
}

type explorerStreamV3 struct {
	timeSourceInitialized bool
	timeDiffSamples       []int64
}

func (g *ExplorerProtoImplV3) readStream(subCtx context.Context, readyChan chan struct{}, initFlag *int32) {
	defer close(g.streamDone)
	defer g.cancelFn()
	DATA_PACKET_HEADER := []byte{0x01, 0xFE}
	DATA_PACKET_TAILER := []byte{0xEF, 0x10}
	packetDelim := append(DATA_PACKET_TAILER, DATA_PACKET_HEADER...)

	for {
		select {
		case <-subCtx.Done():
			g.Logger.Infoln("exiting from data packet reader")
			if atomic.LoadInt32(initFlag) == 0 {
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
			func(bufPtr *[]byte, updatedAt *time.Time) bool {
				buf := *bufPtr
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
			return
		}
		if subCtx.Err() != nil {
			return
		}
		if timeout {
			g.Logger.Errorln("timeout when reading data from transport")
			continue
		}

		if len(recvBuf) == 0 {
			continue
		}
		if len(recvBuf) < 17+len(packetDelim) || !bytes.HasSuffix(recvBuf, packetDelim) {
			g.deviceStatus.IncrementErrors()
			continue
		}
		packetBytes := recvBuf[:len(recvBuf)-len(packetDelim)] // without header and tailer

		g.processPacket(subCtx, packetBytes, recvEndMonotonicTime, recvEndTime, readyChan, initFlag)
	}
}

// Serialize reader-side resets with decoding, including FIFO consumption.
// Otherwise a packet already dequeued by the decoder can repopulate a reset archive.
func (g *ExplorerProtoImplV3) processPacket(subCtx context.Context, packetBytes []byte, recvEndMonotonicTime, recvEndTime time.Time, readyChan chan struct{}, initFlag *int32) {
	g.streamMutex.Lock()
	defer g.streamMutex.Unlock()
	if subCtx.Err() != nil {
		return
	}
	if len(packetBytes) < 17 {
		g.deviceStatus.IncrementErrors()
		return
	}
	variableAllSet := g.variablesReady()

	if err := g.verifyChecksum(packetBytes); err == nil {
		mcuTimestamp := int64(binary.LittleEndian.Uint64(packetBytes[:8]))
		deviceConfig := binary.LittleEndian.Uint32(packetBytes[8 : 8+4])
		gnssEnabled := g.parseGnssAvailibility(deviceConfig)
		channelChunkLength, channelSize, _ := g.getChannelSize(deviceConfig)
		if len(packetBytes) != 17+channelSize {
			g.deviceStatus.IncrementErrors()
			return
		}
		packetInterval := g.parsePacketInterval(deviceConfig)
		sampleRate := g.parseSampleRate(deviceConfig)
		packetLatency := packetInterval + time.Duration(1000/sampleRate*channelChunkLength)*time.Millisecond

		if !g.isDataStreamStable && !gnssEnabled {
			timeDiff := recvEndTime.UnixMilli() - mcuTimestamp - packetLatency.Milliseconds()
			g.stream.timeDiffSamples = append(g.stream.timeDiffSamples, timeDiff)
			if len(g.stream.timeDiffSamples) > STABLE_CHECK_SAMPLES {
				g.stream.timeDiffSamples = g.stream.timeDiffSamples[1:]
			}

			if len(g.stream.timeDiffSamples) == STABLE_CHECK_SAMPLES {
				if minVal, maxVal := lo.Min(g.stream.timeDiffSamples), lo.Max(g.stream.timeDiffSamples); math.Abs(float64(maxVal-minVal)) < 5 {
					g.isDataStreamStable = true
					g.fifoBuffer.Reset()
					g.Logger.Infof("data time series stabilized, final time difference = %d ms", timeDiff)
				} else if (mcuTimestamp/1000)%10 == 0 {
					g.Logger.Warnf("waiting for data time series to settle down, this may take a while, current time difference = %d ms", timeDiff)
				}
			} else if (mcuTimestamp/1000)%2 == 0 {
				g.Logger.Warnln("collecting data time series, this may take a while")
			}
		} else if gnssEnabled {
			g.isDataStreamStable = true
		}

		variableAllSet := g.variablesReady()

		if variableAllSet {
			if gnssEnabled && !g.stream.timeSourceInitialized {
				g.TimeSource.Update(recvEndMonotonicTime, time.UnixMilli(mcuTimestamp).Add(packetLatency), 0, timesource.MonotonicNow)

				g.isDataStreamStable = false
				g.stream.timeSourceInitialized = true
				g.resetFlags()

				g.Logger.Infof("time synchronized with Explorer built-in GNSS module")
			} else if !g.stream.timeSourceInitialized {
				g.Logger.Infoln("synchronizing time with NTP servers, it may take a while")
				offset, err := g.ntpClient.QueryAverageContext(subCtx, NTP_MEASUREMENT_ATTEMPTS)
				if subCtx.Err() != nil {
					return
				}
				if err != nil {
					g.Logger.Errorf("failed to synchronize time with NTP server: %v", err)
					if atomic.LoadInt32(initFlag) == 0 {
						g.cancelFn()
					}
					return
				} else {
					g.Logger.Infof("time synchronized with NTP server, local monotonic time offset: %d ms", offset.Milliseconds())
				}

				currentMonotonicTime := timesource.MonotonicNow()
				g.TimeSource.Update(currentMonotonicTime, currentMonotonicTime.Add(offset), 0, timesource.MonotonicNow)
				g.isDataStreamStable = false
				g.stream.timeSourceInitialized = true
				g.resetFlags()
			}

			if atomic.LoadInt32(initFlag) == 0 {
				atomic.StoreInt32(initFlag, 1)
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
			g.stream.timeDiffSamples = make([]int64, 0, STABLE_CHECK_SAMPLES)
		} else {
			variableAllSet := g.variablesReady()

			if gnssEnabled && g.isDataStreamStable && variableAllSet {
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
	} else if g.isDataStreamStable && variableAllSet {
		g.Logger.Errorln(err)
		g.deviceStatus.IncrementErrors()
	}
}

func (g *ExplorerProtoImplV3) decodeStream(subCtx context.Context, initFlag *int32) {
	defer close(g.decoderDone)
	const decodeInterval = 10 * time.Millisecond
	for timer := time.NewTimer(decodeInterval); ; {
		timer.Reset(decodeInterval)

		select {
		case <-timer.C:
			g.decodePacket(initFlag)
		case <-subCtx.Done():
			g.Logger.Infoln("exiting from data packet decoder")
			timer.Stop()
			return
		}
	}
}

func (g *ExplorerProtoImplV3) decodePacket(initFlag *int32) {
	g.streamMutex.Lock()
	defer g.streamMutex.Unlock()
	dataPackets, err := g.fifoBuffer.Read(1)
	if err != nil {
		return
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

	variableAllSet := g.variablesReady()
	if !variableAllSet {
		if (mcuTimestamp/1000)%4 == 0 {
			g.Logger.Warnln("waiting for device config to be fully collected, this may take a while")
		}
		return
	}

	if atomic.LoadInt32(initFlag) == 0 {
		g.Logger.Warnln("waiting for time to be synchronized, this may take a while")
		return
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
		channelCodes[idx] = ch.ChannelCode
		g.channelDataBuf[idx].ByteSize = ch.ByteSize
		g.channelDataBuf[idx].ChannelCode = ch.ChannelCode
		g.channelDataBuf[idx].ChannelId = ch.ChannelId
		g.channelDataBuf[idx].DataType = ch.DataType
		g.channelDataBuf[idx].Data = append(g.channelDataBuf[idx].Data, ch.Data...)
	}
	g.deviceConfig.SetChannelCodes(channelCodes)

	sampleRate := g.deviceConfig.GetSampleRate()
	g.messageBusRealtime.Publish(NewEvent(timeObj, &g.deviceConfig, lo.Map(
		channelData,
		func(ch *ChannelData, _ int) ChannelData { return *ch },
	)))

	g.flagMutex.Lock()
	collectedSamples := g.collectedSamples
	g.flagMutex.Unlock()

	if collectedSamples < sampleRate {
		return
	} else if collectedSamples == sampleRate {
		g.messageBus.Publish(NewEvent(g.packetTimeObj, &g.deviceConfig, g.channelDataBuf))
		g.deviceStatus.IncrementMessages()
	} else {
		g.Logger.Warnln("collected samples exceeded the sample rate, resetting counters")
		g.resetVariables()
	}

	g.resetFlags()
}

func (g *ExplorerProtoImplV3) synchronizeNTP(subCtx context.Context, readyChan <-chan struct{}) {
	const resyncInterval = NTP_RESYNC_INTERVAL
	defer close(g.ntpDone)
	select {
	case <-readyChan:
	case <-subCtx.Done():
		return
	}

	var prevCalibTime time.Time
	timer := time.NewTimer(resyncInterval)
	defer timer.Stop()
	for {
		select {
		case calibTimeData := <-g.timeCalibrationChan:
			if prevCalibTime.Unix() == calibTimeData[1].Unix() {
				continue
			}
			prevCalibTime = calibTimeData[1]
			g.TimeSource.Update(calibTimeData[0], calibTimeData[1], 0, nil)
		case <-timer.C:
			variableAllSet := g.variablesReady()
			if deviceConfig := g.GetConfig(); deviceConfig.GetGnssAvailability() || !variableAllSet {
				timer.Reset(resyncInterval)
				continue
			}
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
