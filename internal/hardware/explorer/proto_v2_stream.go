package explorer

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/anyshake/observer/pkg/timesource"
	"github.com/samber/lo"
)

const (
	V2_PACKET_SIZE  = 75
	V2_CHANNEL_SIZE = 5
	V2_TIME_JUMP_MS = 5000
	V2_JITTER_MS    = 5
)

// All fields are owned by the packet reader, under streamMutex. The NTP worker
// only reads the mode/generation and updates the cached NTP measurement.
type explorerStreamV2 struct {
	generation       uint64
	modeKnown        bool
	clockInitialized bool
	hasPrevious      bool
	previous         int64
	rateStart        int64
	ratePackets      int
	rateMinDelta     int64
	rateMaxDelta     int64
	timeOffset       int64
	offsetSamples    []int64
	archiveStart     time.Time
	archiveSamples   int
	lastCalibration  int64
	ntpOffset        time.Duration
	ntpMeasuredAt    time.Time
}

// The caller holds streamMutex. A discontinuity invalidates both the time
// domain and any partial archive; the next DEVICE_ID must confirm the mode.
func (g *ExplorerProtoImplV2) resetStream() {
	previous := g.stream
	g.stream = explorerStreamV2{
		generation:    previous.generation + 1,
		ntpOffset:     previous.ntpOffset,
		ntpMeasuredAt: previous.ntpMeasuredAt,
	}
	g.resetVariables()
	g.deviceConfig.SetGnssAvailability(false)
	g.deviceConfig.SetSampleRate(0)
	g.deviceConfig.SetPacketInterval(0)
	g.channelDataBuf = nil
}

func (g *ExplorerProtoImplV2) readStream(ctx context.Context, ready chan<- struct{}) {
	defer close(g.streamDone)
	defer g.cancelFn()
	header := []byte{0xFA, 0xDE}
	buffer := make([]byte, 0, V2_PACKET_SIZE*3)
	chunk := make([]byte, V2_PACKET_SIZE*2)
	started := false
	for ctx.Err() == nil {
		n, readErr := g.Transport.Read(chunk)
		received := timesource.MonotonicNow()
		if ctx.Err() != nil {
			return
		}
		buffer = append(buffer, chunk[:n]...)
		for len(buffer) >= len(header) {
			index := bytes.Index(buffer, header)
			if index < 0 {
				// Preserve a possible header split across two transport reads.
				buffer = buffer[len(buffer)-1:]
				break
			}
			buffer = buffer[index:]
			if len(buffer) < V2_PACKET_SIZE {
				break
			}
			packet := buffer[:V2_PACKET_SIZE]
			if err := g.verifyChecksum(packet, header); err != nil {
				g.deviceStatus.IncrementErrors()
				// A false header must not consume the following valid packet.
				buffer = buffer[1:]
				continue
			}
			initialized, err := g.processPacket(ctx, packet, received, g.Transport.GetLatency(V2_PACKET_SIZE))
			if err != nil {
				g.Logger.Errorf("failed to process v2 packet: %v", err)
				return
			}
			if initialized && !started {
				g.deviceStatus.SetStartedAt(g.TimeSource.Now())
				close(ready)
				started = true
			}
			buffer = buffer[V2_PACKET_SIZE:]
		}
		if readErr != nil {
			g.Logger.Errorf("failed to read data from transport: %v", readErr)
			return
		}
	}
}

func (g *ExplorerProtoImplV2) processPacket(ctx context.Context, packet []byte, received time.Time, latency time.Duration) (bool, error) {
	if len(packet) != V2_PACKET_SIZE || !bytes.Equal(packet[:2], []byte{0xFA, 0xDE}) {
		return false, fmt.Errorf("invalid v2 packet length or header")
	}
	if err := g.verifyChecksum(packet, packet[:2]); err != nil {
		return false, err
	}
	g.streamMutex.Lock()
	defer g.streamMutex.Unlock()
	if err := ctx.Err(); err != nil {
		return false, err
	}
	s := &g.stream
	rawTimestamp := binary.LittleEndian.Uint64(packet[2:10])
	// Bound arithmetic to a supported civil-time range before converting uint64 or
	// adding an offset. Corrupt timestamps must never wrap into valid samples.
	if rawTimestamp > 253402300799999 { // 9999-12-31T23:59:59.999Z
		g.resetStream()
		g.deviceStatus.IncrementErrors()
		return false, nil
	}
	timestamp := int64(rawTimestamp)
	variable := binary.LittleEndian.Uint32(packet[10:14])
	deviceSlot := (timestamp/1000)%4 == 0
	gnss := variable&0x80000000 != 0
	if s.hasPrevious && (timestamp < s.previous || timestamp-s.previous >= V2_TIME_JUMP_MS) {
		g.resetStream()
		g.deviceStatus.IncrementErrors()
	}
	if deviceSlot && s.modeKnown && gnss != g.deviceConfig.GetGnssAvailability() {
		g.resetStream()
	}
	if !s.modeKnown {
		if !deviceSlot {
			return false, nil
		}
		s.modeKnown = true
	}
	// DEVICE_ID is the only authoritative mode indicator in v2. Metadata from
	// other slots is ignored until it arrives, including after a clock jump.
	g.getVariableData(timestamp, variable)
	gnss = g.deviceConfig.GetGnssAvailability()
	if gnss && (timestamp < 946684800000 || timestamp >= 4102444800000) {
		// Explorer GNSS time must be a Unix date in 2000--2099, not MCU uptime.
		g.resetStream()
		g.deviceStatus.IncrementErrors()
		return false, nil
	}
	if !s.hasPrevious {
		s.previous, s.hasPrevious = timestamp, true
		s.rateStart = timestamp
		return false, nil
	}
	delta := timestamp - s.previous
	if delta == 0 {
		return false, nil // Duplicate packets neither add samples nor alter the rate.
	}
	s.previous = timestamp
	interval := g.deviceConfig.GetPacketInterval().Milliseconds()
	jitter := min(int64(V2_JITTER_MS), interval/4)
	if delta > 1000 || (interval > 0 && (delta < interval-jitter || delta > interval+jitter)) {
		// Missing packets or a rate change cannot be appended to a contiguous
		// one-second archive. Re-estimate the interval from the next pair.
		g.deviceConfig.SetSampleRate(0)
		g.deviceConfig.SetPacketInterval(0)
		g.channelDataBuf = nil
		s.archiveSamples = 0
		s.rateStart, s.ratePackets = timestamp, 0
		g.deviceStatus.IncrementErrors()
		return false, nil
	}
	if interval == 0 {
		// Infer the rate over a full second, as in the original v2 decoder.
		// One interval with 1ms jitter must not turn 125Hz into 122Hz.
		s.ratePackets++
		if s.ratePackets == 1 {
			s.rateMinDelta, s.rateMaxDelta = delta, delta
		} else {
			s.rateMinDelta = min(s.rateMinDelta, delta)
			s.rateMaxDelta = max(s.rateMaxDelta, delta)
		}
		elapsed := timestamp - s.rateStart
		if elapsed < 1000-V2_JITTER_MS {
			return false, nil
		}
		estimatedInterval := int64(1000 / s.ratePackets)
		jitter := min(int64(V2_JITTER_MS), estimatedInterval/4)
		if elapsed > 1000+V2_JITTER_MS || s.rateMinDelta < estimatedInterval-jitter || s.rateMaxDelta > estimatedInterval+jitter {
			s.rateStart, s.ratePackets = timestamp, 0
			return false, nil
		}
		sampleRate := s.ratePackets * V2_CHANNEL_SIZE
		g.deviceConfig.SetSampleRate(sampleRate)
		g.deviceConfig.SetPacketInterval(time.Second * V2_CHANNEL_SIZE / time.Duration(sampleRate))
	}
	if !g.variableAllSet {
		return false, nil
	}
	packetDelay := g.deviceConfig.GetPacketInterval() + latency
	if !s.clockInitialized {
		if gnss {
			g.TimeSource.Update(received, time.UnixMilli(timestamp).Add(packetDelay), 0, timesource.MonotonicNow)
		} else {
			if s.ntpMeasuredAt.IsZero() || timesource.MonotonicNow().Sub(s.ntpMeasuredAt) >= NTP_RESYNC_INTERVAL {
				offset, err := g.ntpClient.QueryAverageContext(ctx, NTP_MEASUREMENT_ATTEMPTS)
				if err != nil {
					return false, fmt.Errorf("failed to synchronize time with NTP: %w", err)
				}
				s.ntpOffset, s.ntpMeasuredAt = offset, timesource.MonotonicNow()
			}
			if err := ctx.Err(); err != nil {
				return false, err
			}
			now := timesource.MonotonicNow()
			g.TimeSource.Update(now, now.Add(s.ntpOffset), 0, timesource.MonotonicNow)
			s.timeOffset = g.TimeSource.Now().Add(received.Sub(timesource.MonotonicNow())-packetDelay).UnixMilli() - timestamp
		}
		s.clockInitialized = true
	}
	packetTime := time.UnixMilli(timestamp)
	if gnss {
		if timestamp/1000 != s.lastCalibration {
			g.TimeSource.Update(received, packetTime.Add(packetDelay), 0, timesource.MonotonicNow)
			s.lastCalibration = timestamp / 1000
		}
	} else {
		// Continuously refresh the drift window. Only use it while receive
		// timing is stable; a burst of buffered packets is not a clock offset.
		offset := g.TimeSource.Now().Add(received.Sub(timesource.MonotonicNow())-packetDelay).UnixMilli() - timestamp
		s.offsetSamples = append(s.offsetSamples, offset)
		if len(s.offsetSamples) > STABLE_CHECK_SAMPLES {
			s.offsetSamples = s.offsetSamples[1:]
		}
		if len(s.offsetSamples) == STABLE_CHECK_SAMPLES && lo.Max(s.offsetSamples)-lo.Min(s.offsetSamples) <= V2_JITTER_MS {
			s.timeOffset = lo.Mean(s.offsetSamples)
		}
		packetTime = time.UnixMilli(timestamp + s.timeOffset)
	}
	channelData, err := g.getChannelData(packet, 2, V2_CHANNEL_SIZE)
	if err != nil {
		return false, err
	}
	if s.archiveSamples == 0 {
		s.archiveStart = packetTime
	}
	s.archiveSamples += V2_CHANNEL_SIZE
	g.messageBusRealtime.Publish(NewEvent(packetTime, &g.deviceConfig, channelData))
	if s.archiveSamples >= g.deviceConfig.GetSampleRate() {
		if s.archiveSamples == g.deviceConfig.GetSampleRate() {
			g.messageBus.Publish(NewEvent(s.archiveStart, &g.deviceConfig, g.channelDataBuf))
			g.deviceStatus.IncrementMessages()
		}
		g.channelDataBuf = nil
		s.archiveSamples = 0
	}
	g.deviceStatus.IncrementFrames()
	g.deviceStatus.SetUpdatedAt(packetTime)
	return true, nil
}

func (g *ExplorerProtoImplV2) synchronizeNTP(ctx context.Context) {
	defer close(g.ntpDone)
	timer := time.NewTicker(NTP_RESYNC_INTERVAL)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			g.streamMutex.Lock()
			generation := g.stream.generation
			enabled := g.stream.clockInitialized && !g.deviceConfig.GetGnssAvailability()
			g.streamMutex.Unlock()
			if !enabled {
				continue
			}
			offset, server, err := g.ntpClient.QueryContext(ctx)
			if ctx.Err() != nil {
				return
			}
			if err != nil {
				g.Logger.Warnf("error occurred while re-synchronizing time with NTP: %v", err)
				continue
			}
			if g.applyNTPMeasurement(generation, offset) {
				g.Logger.Infof("time synchronized with NTP server: %s, local monotonic time offset: %d ms", server, offset.Milliseconds())
			}
		}
	}
}

func (g *ExplorerProtoImplV2) applyNTPMeasurement(generation uint64, offset time.Duration) bool {
	g.streamMutex.Lock()
	defer g.streamMutex.Unlock()
	// A query started in a previous time domain must not overwrite its successor,
	// even if GNSS has since switched back to NTP.
	if g.stream.generation != generation || !g.stream.clockInitialized || g.deviceConfig.GetGnssAvailability() {
		return false
	}
	now := timesource.MonotonicNow()
	g.stream.ntpOffset, g.stream.ntpMeasuredAt = offset, now
	g.clockDriftBuf.Push(clockDrift{offset: offset, measuredAt: now})
	ppm := getLongTermClockDriftPPM(g.clockDriftBuf, NTP_PPM_MEASURE_WINDOW)
	g.TimeSource.Update(now, now.Add(offset), ppm, timesource.MonotonicNow)
	g.stream.offsetSamples = nil
	return true
}
