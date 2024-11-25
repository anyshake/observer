package explorer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"
	"unsafe"

	"github.com/alphadose/haxmap"
	"github.com/anyshake/observer/utils/fifo"
	"github.com/sirupsen/logrus"
	messagebus "github.com/vardius/message-bus"
)

var (
	LEGACY_PACKET_FRAME_HEADER   = []byte{0xFC, 0x1B}
	MAINLINE_PACKET_FRAME_HEADER = []byte{0xFA, 0xDE}
)

// In legacy mode, each packet contains 3 channels, n samples per channel.
// The packet is sent at an interval of (1000 / sample rate) milliseconds.
// Set n = 5 (also in Explorer) fits the common sample rates (25, 50, 100, 125 Hz).
const LEGACY_PACKET_CHANNEL_SIZE = 5

// Legacy packet structure, fixed size.
// Each channel has a checksum, which is the XOR of all bytes in the channel.
type legacyPacket struct {
	Z_Axis   [LEGACY_PACKET_CHANNEL_SIZE]int32
	E_Axis   [LEGACY_PACKET_CHANNEL_SIZE]int32
	N_Axis   [LEGACY_PACKET_CHANNEL_SIZE]int32
	Checksum [3]uint8
}

func (g *legacyPacket) length() int {
	return int(unsafe.Sizeof(*g))
}

func (g *legacyPacket) decode(data []byte) error {
	err := binary.Read(bytes.NewReader(data), binary.LittleEndian, g)
	if err != nil {
		return err
	}

	// Using XOR algorithm
	calcChecksum := [3]uint8{0, 0, 0}
	zAxisOffset := int(unsafe.Sizeof(g.Z_Axis))
	for i := 0; i < zAxisOffset; i++ {
		calcChecksum[0] ^= data[i]
	}
	eAxisOffset := zAxisOffset + int(unsafe.Sizeof(g.E_Axis))
	for i := zAxisOffset; i < eAxisOffset; i++ {
		calcChecksum[1] ^= data[i]
	}
	nAxisOffset := eAxisOffset + int(unsafe.Sizeof(g.N_Axis))
	for i := eAxisOffset; i < nAxisOffset; i++ {
		calcChecksum[2] ^= data[i]
	}
	for i := 0; i < len(calcChecksum); i++ {
		if calcChecksum[i] != g.Checksum[i] {
			return fmt.Errorf("checksum mismatch, expected %v, got %v", g.Checksum, calcChecksum)
		}
	}

	return nil
}

// In mainline mode, each packet contains 3 channels, n samples per channel.
// The packet is sent at an interval of (1000 / sample rate) milliseconds.
// Set n = 5 (also in Explorer) fits the common sample rates (25, 50, 100, 125 Hz).
const MAINLINE_PACKET_CHANNEL_SIZE = 5

// Mainline packet header structure, fixed size.
// The VariableData be Device ID, Latitude, Longitude, Elevation in int32 / float32 format.
type mainlinePacket struct {
	Timestamp    int64
	VariableData [4]byte // Can be int32 or float32
	VariableName string  // Exclude from length calculation
	Z_axis       [MAINLINE_PACKET_CHANNEL_SIZE]int32
	E_axis       [MAINLINE_PACKET_CHANNEL_SIZE]int32
	N_axis       [MAINLINE_PACKET_CHANNEL_SIZE]int32
	Checksum     uint8
}

func (g *mainlinePacket) length() int {
	return int(unsafe.Sizeof(g.Timestamp) +
		unsafe.Sizeof(g.VariableData) +
		unsafe.Sizeof(g.Z_axis) +
		unsafe.Sizeof(g.E_axis) +
		unsafe.Sizeof(g.N_axis) +
		unsafe.Sizeof(g.Checksum))
}

func (g *mainlinePacket) decode(data []byte) error {
	// Restore header checksum, note that the byte order is little-endian
	checksumIndex := len(data) - int(unsafe.Sizeof(g.Checksum))
	g.Checksum = data[checksumIndex]

	// Using XOR algorithm to calculate the header checksum
	calcHeaderChecksum := uint8(0)
	for i := 0; i < checksumIndex; i++ {
		calcHeaderChecksum ^= data[i]
	}
	if calcHeaderChecksum != g.Checksum {
		return fmt.Errorf("header checksum mismatch, expected %d, got %d", g.Checksum, calcHeaderChecksum)
	}

	// Restore the header data, note that the byte order is little-endian
	g.Timestamp = int64(binary.LittleEndian.Uint64(data[:unsafe.Sizeof(g.Timestamp)]))
	switch (g.Timestamp / time.Second.Milliseconds()) % 4 {
	case 0:
		g.VariableName = "device_info"
	case 1:
		g.VariableName = "latitude"
	case 2:
		g.VariableName = "longitude"
	case 3:
		g.VariableName = "elevation"
	}
	variableDataIndex := int(unsafe.Sizeof(g.Timestamp) + unsafe.Sizeof(g.VariableData))
	copy(g.VariableData[:], data[unsafe.Sizeof(g.Timestamp):variableDataIndex])

	// Restore the channel data, note that the byte order is little-endian
	zAxisOffset := variableDataIndex + int(unsafe.Sizeof(g.Z_axis))
	err := binary.Read(bytes.NewReader(data[variableDataIndex:zAxisOffset]), binary.LittleEndian, g.Z_axis[:])
	if err != nil {
		return err
	}
	eAxisOffset := zAxisOffset + int(unsafe.Sizeof(g.E_axis))
	err = binary.Read(bytes.NewReader(data[zAxisOffset:eAxisOffset]), binary.LittleEndian, g.E_axis[:])
	if err != nil {
		return err
	}
	nAxisOffset := eAxisOffset + int(unsafe.Sizeof(g.N_axis))
	return binary.Read(bytes.NewReader(data[eAxisOffset:nAxisOffset]), binary.LittleEndian, g.N_axis[:])
}

type ExplorerDriverImpl struct {
	logger         *logrus.Entry
	legacyPacket   legacyPacket
	mainlinePacket mainlinePacket
}

func (e *ExplorerDriverImpl) handleReadLegacyPacket(deps *ExplorerDependency, fifoBuffer *fifo.Buffer) {
	// Set to 0xFFFFFFFF to indicate legacy mode
	deps.Config.SetDeviceInfo(math.MaxUint32)

	findIndices := func(arr []byte, sep []byte) []int {
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

	// Read data from the transport continuously
	go func() {
		var (
			recvSize        = len(LEGACY_PACKET_FRAME_HEADER) + e.legacyPacket.length()
			timeBytes       = make([]byte, 8)                       // 8 bytes for int64 timestamp
			recvBuf         = make([]byte, recvSize)                // Received buffer from transport
			packetBuf       = make([]byte, recvSize+len(timeBytes)) // Packet buffer with timestamp
			prevHeaderIndex = -1                                    // Last index of the header
		)
		for {
			select {
			case <-deps.CancelToken.Done():
				e.logger.Infof("cancelling read data from transport")
				return
			default:
				_, err := deps.Transport.Read(recvBuf, time.Second, false)
				if err != nil {
					e.logger.Errorf("failed to read data from transport: %v", err)
					continue
				}

				// Record the current time of the packet
				currentTime := deps.FallbackTime.Get()
				binary.BigEndian.PutUint64(timeBytes, uint64(currentTime.UnixMilli()))

				// Find possible header in the buffer to insert current time next to the header
				headerIndices := findIndices(recvBuf, LEGACY_PACKET_FRAME_HEADER)
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
				copy(packetBuf, recvBuf[:headerIndex+len(LEGACY_PACKET_FRAME_HEADER)])                                                 // Copy header
				copy(packetBuf[headerIndex+len(LEGACY_PACKET_FRAME_HEADER):headerIndex+len(LEGACY_PACKET_FRAME_HEADER)+8], timeBytes)  // Copy timestamp
				copy(packetBuf[headerIndex+len(LEGACY_PACKET_FRAME_HEADER)+8:], recvBuf[headerIndex+len(LEGACY_PACKET_FRAME_HEADER):]) // Copy packet
				fifoBuffer.Write(packetBuf)
			}
		}
	}()

	findClosestSampleRate := func(currentSampleRate int) int {
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

	// Read data from the FIFO buffer continuously
	var (
		packetBuffer = []mainlinePacket{}                                            // Legacy packet is converted to mainline packet internally
		readSize     = len(LEGACY_PACKET_FRAME_HEADER) + 8 + e.legacyPacket.length() // header + timestamp + legacyPacket
		nextTick     = int64(0)                                                      // Expected timestamp for the next published data on message bus
		timer        = time.NewTimer(time.Millisecond)
	)
	for {
		timer.Reset(time.Millisecond)

		select {
		case <-deps.CancelToken.Done():
			return
		case <-timer.C:
			dat, err := fifoBuffer.Peek(LEGACY_PACKET_FRAME_HEADER, readSize)
			if err != nil {
				continue
			}

			// Read and decode the legacy packet
			err = e.legacyPacket.decode(dat[len(LEGACY_PACKET_FRAME_HEADER)+8:])
			if err != nil {
				e.logger.Errorf("failed to decode legacy packet: %v", err)
				deps.Health.SetErrors(deps.Health.GetErrors() + 1)
				continue
			}

			// Extract timestamp from the buffer
			timestamp := int64(binary.BigEndian.Uint64(
				dat[len(LEGACY_PACKET_FRAME_HEADER) : len(LEGACY_PACKET_FRAME_HEADER)+8],
			))

			// Append the packet to the buffer
			if nextTick == 0 {
				nextTick = timestamp
			} else {
				packetBuffer = append(packetBuffer, mainlinePacket{
					Timestamp: timestamp,
					Z_axis:    e.legacyPacket.Z_Axis,
					E_axis:    e.legacyPacket.E_Axis,
					N_axis:    e.legacyPacket.N_Axis,
				})
			}

			// Calculate proper sample rate to avoid jitter
			currentSampleRate := len(packetBuffer) * MAINLINE_PACKET_CHANNEL_SIZE
			targetSampleRate := findClosestSampleRate(currentSampleRate)

			if math.Abs(float64(timestamp-nextTick)) <= EXPLORER_ALLOWED_JITTER_MS && currentSampleRate == targetSampleRate {
				// Update the next tick even if the buffer is empty
				nextTick = timestamp + time.Second.Milliseconds()
				if len(packetBuffer) == 0 {
					continue
				}

				// Merge the packet buffer into a single packet
				var (
					z_axis_count []int32
					e_axis_count []int32
					n_axis_count []int32
				)
				for _, packet := range packetBuffer {
					z_axis_count = append(z_axis_count, packet.Z_axis[:]...)
					e_axis_count = append(e_axis_count, packet.E_axis[:]...)
					n_axis_count = append(n_axis_count, packet.N_axis[:]...)
				}

				// Publish the final packet
				finalPacket := ExplorerData{
					SampleRate: targetSampleRate,
					Z_Axis:     z_axis_count,
					E_Axis:     e_axis_count,
					N_Axis:     n_axis_count,
					Timestamp:  timestamp - time.Second.Milliseconds(),
				}
				deps.messageBus.Publish("explorer", &finalPacket)

				// Update the health status
				deps.Health.SetSampleRate(targetSampleRate)
				deps.Health.SetReceived(deps.Health.GetReceived() + 1)
				deps.Health.SetUpdatedAt(time.UnixMilli(timestamp).UTC())

				packetBuffer = []mainlinePacket{}
			} else if timestamp-nextTick > EXPLORER_ALLOWED_JITTER_MS {
				// Update the next tick, clear the buffer if the jitter exceeds the threshold
				nextTick = timestamp + time.Second.Milliseconds()
				packetBuffer = []mainlinePacket{}
			}
		}
	}
}

func (e *ExplorerDriverImpl) handleReadMainlinePacket(deps *ExplorerDependency, fifoBuffer *fifo.Buffer) {
	recvSize := len(MAINLINE_PACKET_FRAME_HEADER) + e.mainlinePacket.length()

	// Read data from the transport continuously
	go func() {
		buf := make([]byte, recvSize)
		for {
			select {
			case <-deps.CancelToken.Done():
				e.logger.Infof("cancelling read data from transport")
				return
			default:
				n, err := deps.Transport.Read(buf, time.Second, false)
				if err != nil {
					e.logger.Errorf("failed to read data from transport: %v", err)
					continue
				}

				fifoBuffer.Write(buf[:n])
			}
		}
	}()

	// Read data from the FIFO buffer continuously
	var (
		packetBuffer = []mainlinePacket{}
		noGnssMode   = false
		nextTick     = int64(0) // Expected timestamp for the next published data on message bus
		timeDiff     = int64(0) // For non-GNSS mode, we need time difference between the packet and NTP time
		timer        = time.NewTimer(time.Millisecond)
	)
	for {
		timer.Reset(time.Millisecond)

		select {
		case <-deps.CancelToken.Done():
			return
		case <-timer.C:
			dat, err := fifoBuffer.Peek(MAINLINE_PACKET_FRAME_HEADER, recvSize)
			if err != nil {
				continue
			}
			err = e.mainlinePacket.decode(dat[len(MAINLINE_PACKET_FRAME_HEADER):])
			if err != nil {
				e.logger.Errorf("failed to decode mainline packet: %v", err)
				deps.Health.SetErrors(deps.Health.GetErrors() + 1)
				continue
			}

			// If device ID is not initialized, get the device ID from the packet
			_, deviceId := deps.Config.GetDeviceInfo()
			if deviceId == 0 && e.mainlinePacket.VariableName == "device_info" {
				deviceInfo := binary.LittleEndian.Uint32(e.mainlinePacket.VariableData[:])
				// When the most significant bit is 0, it means the device is running without GNSS module.
				// In this case, the latitude, longitude, elevation will not be updated.
				// The timestamp will be replaced with the NTP time.
				if deviceInfo&0x80000000 == 0 {
					noGnssMode = true
				}
				deps.Config.SetDeviceInfo(deviceInfo)
				_, deviceId = deps.Config.GetDeviceInfo()
				e.logger.Infof("got current device ID: %08X", deviceId)
			} else if deviceId == 0 {
				continue
			}

			// Update the latitude, longitude, elevation
			switch e.mainlinePacket.VariableName {
			case "latitude":
				latitude := math.Float32frombits(binary.LittleEndian.Uint32(e.mainlinePacket.VariableData[:]))
				if latitude >= -90 && latitude <= 90 && !noGnssMode {
					deps.Config.SetLatitude(float64(latitude))
				}
			case "longitude":
				longitude := math.Float32frombits(binary.LittleEndian.Uint32(e.mainlinePacket.VariableData[:]))
				if longitude >= -180 && longitude <= 180 && !noGnssMode {
					deps.Config.SetLongitude(float64(longitude))
				}
			case "elevation":
				elevation := math.Float32frombits(binary.LittleEndian.Uint32(e.mainlinePacket.VariableData[:]))
				if elevation >= 0 && !noGnssMode {
					deps.Config.SetElevation(float64(elevation))
				}
			}

			if noGnssMode {
				// Update time difference at the first packet or everyday at 00:00:00 UTC
				currentTime := deps.FallbackTime.Get()
				if nextTick == 0 || currentTime.Unix()%(int64(time.Hour.Seconds())*24) == 0 {
					timeDiff = currentTime.UTC().UnixMilli() - e.mainlinePacket.Timestamp
				}
			} else {
				timeDiff = 0
			}

			// Append the packet to the buffer
			if nextTick == 0 {
				nextTick = e.mainlinePacket.Timestamp
			} else {
				packetBuffer = append(packetBuffer, e.mainlinePacket)
			}

			if math.Abs(float64(e.mainlinePacket.Timestamp-nextTick)) <= EXPLORER_ALLOWED_JITTER_MS {
				// Update the next tick even if the buffer is empty
				nextTick = e.mainlinePacket.Timestamp + time.Second.Milliseconds()
				if len(packetBuffer) == 0 {
					continue
				}

				// Merge the packet buffer into a single packet
				var (
					z_axis_count []int32
					e_axis_count []int32
					n_axis_count []int32
				)
				for _, packet := range packetBuffer {
					z_axis_count = append(z_axis_count, packet.Z_axis[:]...)
					e_axis_count = append(e_axis_count, packet.E_axis[:]...)
					n_axis_count = append(n_axis_count, packet.N_axis[:]...)
				}

				// Publish the final packet
				sampleRate := len(packetBuffer) * MAINLINE_PACKET_CHANNEL_SIZE
				finalPacket := ExplorerData{
					SampleRate: sampleRate,
					Z_Axis:     z_axis_count,
					E_Axis:     e_axis_count,
					N_Axis:     n_axis_count,
					Timestamp:  (e.mainlinePacket.Timestamp - time.Second.Milliseconds()) + timeDiff,
				}
				deps.messageBus.Publish("explorer", &finalPacket)

				// Update the health status
				deps.Health.SetSampleRate(sampleRate)
				deps.Health.SetReceived(deps.Health.GetReceived() + 1)
				deps.Health.SetUpdatedAt(time.UnixMilli(e.mainlinePacket.Timestamp).UTC())

				packetBuffer = []mainlinePacket{}
			} else if nextTick-e.mainlinePacket.Timestamp > time.Second.Milliseconds()+EXPLORER_ALLOWED_JITTER_MS {
				// Update the next tick, clear the buffer if the jitter exceeds the threshold
				nextTick = e.mainlinePacket.Timestamp + time.Second.Milliseconds()
				packetBuffer = []mainlinePacket{}
				// Update the time difference again if the device is running without GNSS module
				if noGnssMode {
					currentTime := deps.FallbackTime.Get()
					timeDiff = currentTime.UTC().UnixMilli() - e.mainlinePacket.Timestamp
				}
			}
		}
	}
}

func (e *ExplorerDriverImpl) readerDaemon(deps *ExplorerDependency) {
	fifoBuffer := fifo.New(8192)

	if deps.Config.GetLegacyMode() {
		e.handleReadLegacyPacket(deps, &fifoBuffer)
	} else {
		e.handleReadMainlinePacket(deps, &fifoBuffer)
	}
}

func (e *ExplorerDriverImpl) Init(deps *ExplorerDependency, logger *logrus.Entry) error {
	e.logger = logger

	currentTime := deps.FallbackTime.Get()
	deps.Health.SetStartTime(currentTime)
	deps.subscribers = haxmap.New[string, ExplorerEventHandler]()
	deps.messageBus = messagebus.New(1024)

	go e.readerDaemon(deps)
	return nil
}

func (e *ExplorerDriverImpl) Subscribe(deps *ExplorerDependency, clientId string, handler ExplorerEventHandler) error {
	if _, ok := deps.subscribers.Get(clientId); ok {
		return errors.New("this client has already subscribed")
	}
	deps.subscribers.Set(clientId, handler)
	deps.messageBus.Subscribe("explorer", handler)
	return nil
}

func (e *ExplorerDriverImpl) Unsubscribe(deps *ExplorerDependency, clientId string) error {
	fn, ok := deps.subscribers.Get(clientId)
	if !ok {
		return errors.New("this client has not subscribed")
	}

	deps.subscribers.Del(clientId)
	return deps.messageBus.Unsubscribe("explorer", fn)
}
