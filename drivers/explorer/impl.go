package explorer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"
	"unsafe"

	"github.com/anyshake/observer/utils/fifo"
	cmap "github.com/orcaman/concurrent-map/v2"
	messagebus "github.com/vardius/message-bus"
)

var (
	legacy_packet_frame_header   = []byte{0xFC, 0x1B}
	mainline_packet_frame_header = []byte{0xF1, 0xD9}
	mainline_packet_frame_tail   = []byte{0xD9, 0xF1}
)

// In legacy mode, each packet contains 3 channels, n samples per channel.
// The packet is sent at an interval of (1000 / sample rate) milliseconds.
// Set n = 5 (also in Explorer) fits the common sample rates (25, 50, 100, 125 Hz).
const legacy_packet_channel_size = 5

// Legacy packet structure, fixed size.
// Each channel has a checksum, which is the XOR of all bytes in the channel.
type legacyPacket struct {
	Z_Axis   [legacy_packet_channel_size]int32
	E_Axis   [legacy_packet_channel_size]int32
	N_Axis   [legacy_packet_channel_size]int32
	Checksum [3]uint8
}

func (g *legacyPacket) length() int {
	return int(unsafe.Sizeof(g.Z_Axis) + unsafe.Sizeof(g.E_Axis) + unsafe.Sizeof(g.N_Axis) + unsafe.Sizeof(g.Checksum))
}

func (g *legacyPacket) decode(data []byte) error {
	err := binary.Read(bytes.NewReader(data), binary.LittleEndian, g)
	if err != nil {
		return err
	}

	// Using XOR algorithm
	calc_checksum := [3]uint8{0, 0, 0}
	z_axis_offset := int(unsafe.Sizeof(g.Z_Axis))
	for i := 0; i < z_axis_offset; i++ {
		calc_checksum[0] ^= data[i]
	}
	e_axis_offset := z_axis_offset + int(unsafe.Sizeof(g.E_Axis))
	for i := z_axis_offset; i < e_axis_offset; i++ {
		calc_checksum[1] ^= data[i]
	}
	n_axis_offset := e_axis_offset + int(unsafe.Sizeof(g.N_Axis))
	for i := e_axis_offset; i < n_axis_offset; i++ {
		calc_checksum[2] ^= data[i]
	}
	for i := 0; i < len(calc_checksum); i++ {
		if calc_checksum[i] != g.Checksum[i] {
			return fmt.Errorf("checksum mismatch, expected %v, got %v", g.Checksum, calc_checksum)
		}
	}

	return nil
}

// Mainline packet header structure, fixed size.
// 34 bytes of header data without the frame header bytes.
type mainlinePacketHeader struct {
	sampleRate uint16
	timestamp  int64
	deviceId   uint32
	latitude   float32
	longitude  float32
	elevation  float32
	reserved   uint64
	checksum   uint8
}

func (g *mainlinePacketHeader) length() int {
	return int(unsafe.Sizeof(g.sampleRate) +
		unsafe.Sizeof(g.timestamp) +
		unsafe.Sizeof(g.deviceId) +
		unsafe.Sizeof(g.latitude) +
		unsafe.Sizeof(g.longitude) +
		unsafe.Sizeof(g.elevation) +
		unsafe.Sizeof(g.reserved) +
		unsafe.Sizeof(g.checksum))
}

func (g *mainlinePacketHeader) decode(data []byte) error {
	g.checksum = data[len(data)-1]

	// Using XOR algorithm
	calc_checksum := uint8(0)
	for i := 0; i < len(data[:34]); i++ {
		calc_checksum ^= data[i]
	}
	if calc_checksum != g.checksum {
		return fmt.Errorf("checksum mismatch, expected %d, got %d", g.checksum, calc_checksum)
	}

	g.sampleRate = binary.LittleEndian.Uint16(data[:2])
	g.timestamp = int64(binary.LittleEndian.Uint64(data[2:10]))
	g.deviceId = binary.LittleEndian.Uint32(data[10:14])
	g.latitude = math.Float32frombits(binary.LittleEndian.Uint32(data[14:18]))
	g.longitude = math.Float32frombits(binary.LittleEndian.Uint32(data[18:22]))
	g.elevation = math.Float32frombits(binary.LittleEndian.Uint32(data[22:26]))
	g.reserved = binary.LittleEndian.Uint64(data[26:34])

	return nil
}

// Mainline packet channel structure, variable number of samples.
// Flexibly sized packet channel depending on the sample rate.
type mainlinePacketChannel struct {
	z_axis   []int32
	e_axis   []int32
	n_axis   []int32
	checksum uint32
}

func (g *mainlinePacketChannel) length(sampleRate int) int {
	return 3*sampleRate*int(unsafe.Sizeof(int32(0))) + // Z, E, N axis data
		int(unsafe.Sizeof(g.checksum)) // Checksum of Z, E, N axis
}

func (g *mainlinePacketChannel) decode(data []byte, sampleRate int) error {
	g.checksum = binary.LittleEndian.Uint32(data[len(data)-4:])

	// Convert little-endian to big-endian for checksum calculation
	for i := 0; i < len(data)-4; i += 4 {
		data[i], data[i+1], data[i+2], data[i+3] = data[i+3], data[i+2], data[i+1], data[i]
	}

	// Using CRC-32/MPEG-2 algorithm
	calc_checksum := uint32(0xFFFFFFFF)
	for _, v := range data[:len(data)-4] {
		calc_checksum ^= uint32(v) << 24
		for i := 0; i < 8; i++ {
			if (calc_checksum & 0x80000000) != 0 {
				calc_checksum = (calc_checksum << 1) ^ 0x04C11DB7
			} else {
				calc_checksum <<= 1
			}
		}
	}
	if calc_checksum != g.checksum {
		return fmt.Errorf("checksum mismatch, expected %d, got %d", g.checksum, calc_checksum)
	}

	// Restore the original data, note that the byte order is big-endian
	g.z_axis = make([]int32, sampleRate)
	binary.Read(bytes.NewReader(data[:sampleRate*int(unsafe.Sizeof(int32(0)))]), binary.BigEndian, g.z_axis)
	g.e_axis = make([]int32, sampleRate)
	binary.Read(bytes.NewReader(data[sampleRate*int(unsafe.Sizeof(int32(0))):2*sampleRate*int(unsafe.Sizeof(int32(0)))]), binary.BigEndian, g.e_axis)
	g.n_axis = make([]int32, sampleRate)
	binary.Read(bytes.NewReader(data[2*sampleRate*int(unsafe.Sizeof(int32(0))):3*sampleRate*int(unsafe.Sizeof(int32(0)))]), binary.BigEndian, g.n_axis)

	return nil
}

// Mainline packet tail structure, fixed size
// 9 bytes of tail data without the frame tail bytes
type mainlinePacketTail struct {
	reserved uint64
	checksum uint8
}

func (g *mainlinePacketTail) length() int {
	return int(unsafe.Sizeof(g.reserved) + unsafe.Sizeof(g.checksum))
}

func (g *mainlinePacketTail) decode(data []byte) error {
	g.checksum = data[8]

	// Using XOR algorithm
	calc_checksum := uint8(0)
	for i := 0; i < len(data); i++ {
		calc_checksum ^= data[i]
	}
	if calc_checksum != g.checksum {
		return fmt.Errorf("checksum mismatch, expected %d, got %d", g.checksum, calc_checksum)
	}

	g.reserved = binary.LittleEndian.Uint64(data[:8])
	return nil
}

type ExplorerDriverImpl struct {
	// Dependencies for legacy mode
	legacyPacket legacyPacket
	// Dependencies for mainline mode
	mainlinePacketHeader  mainlinePacketHeader
	mainlinePacketChannel mainlinePacketChannel
	mainlinePacketTail    mainlinePacketTail
}

func (e *ExplorerDriverImpl) handleReadLegacyPacket(deps *ExplorerDependency) {
	fifoBuffer := fifo.New(e.legacyPacket.length() * 512)

	// Read data from the transport continuously
	go func() {
		buf := make([]byte, e.legacyPacket.length())
		for {
			select {
			case <-deps.CancelToken.Done():
				return
			default:
				n, err := deps.Transport.Read(buf, 10*time.Millisecond, false)
				if err != nil {
					return
				}

				fifoBuffer.Write(buf[:n])
			}
		}
	}()

	// Reference: https://stackoverflow.com/a/51424566
	// Calculate the duration to the next whole second to allivate the drift
	calcDuration := func(currentTime time.Time, duration time.Duration) time.Duration {
		return currentTime.Round(duration).Add(duration).Sub(currentTime)
	}

	// Read data from the FIFO buffer continuously
	var (
		dataBuffer = []legacyPacket{}
		ticker     = time.NewTimer(calcDuration(time.Now(), time.Second))
	)
	for {
		select {
		case <-deps.CancelToken.Done():
			return
		case currentTick := <-ticker.C:
			if len(dataBuffer) > 0 {
				currentTime, err := deps.FallbackTime.Get()
				if err != nil {
					continue
				}

				deps.Health.UpdatedAt = currentTime
				deps.Health.Received++

				var (
					z_axis_count []int32
					e_axis_count []int32
					n_axis_count []int32
				)
				for _, packet := range dataBuffer {
					z_axis_count = append(z_axis_count, packet.Z_Axis[:]...)
					e_axis_count = append(e_axis_count, packet.E_Axis[:]...)
					n_axis_count = append(n_axis_count, packet.N_Axis[:]...)
				}

				sampleRate := len(dataBuffer) * legacy_packet_channel_size
				deps.Health.SampleRate = sampleRate
				finalPacket := ExplorerData{
					SampleRate: sampleRate,
					Z_Axis:     z_axis_count,
					E_Axis:     e_axis_count,
					N_Axis:     n_axis_count,
					Timestamp:  currentTime.UTC().UnixMilli(),
				}
				deps.messageBus.Publish("explorer", &finalPacket)
				dataBuffer = []legacyPacket{}

				ticker.Reset(calcDuration(currentTick, time.Second))
			}
		case <-time.After(500 * time.Microsecond):
			dat, err := fifoBuffer.Read(legacy_packet_frame_header, len(legacy_packet_frame_header)+e.legacyPacket.length())
			if err == nil {
				// Read the packet data
				err = e.legacyPacket.decode(dat[len(legacy_packet_frame_header):])
				if err != nil {
					deps.Health.Errors++
				} else {
					dataBuffer = append(dataBuffer, e.legacyPacket)
				}
			}
		}
	}
}

func (e *ExplorerDriverImpl) handleReadMainlinePacket(deps *ExplorerDependency) {
	for {
		select {
		case <-deps.CancelToken.Done():
			return
		default:
			// Find the header sync bytes
			ok, _ := deps.Transport.Filter(mainline_packet_frame_header, 2*time.Second)
			if !ok {
				continue
			}

			// Read header section and update dependency data
			headerBuf := make([]byte, e.mainlinePacketHeader.length())
			_, err := deps.Transport.Read(headerBuf, time.Second, false)
			if err != nil {
				continue
			}
			err = e.mainlinePacketHeader.decode(headerBuf)
			if err != nil {
				deps.Health.Errors++
				continue
			}
			if e.mainlinePacketHeader.latitude != 0 && e.mainlinePacketHeader.longitude != 0 && e.mainlinePacketHeader.elevation != 0 {
				deps.Config.Latitude = float64(e.mainlinePacketHeader.latitude)
				deps.Config.Longitude = float64(e.mainlinePacketHeader.longitude)
				deps.Config.Elevation = float64(e.mainlinePacketHeader.elevation)
			}

			// Get data section packet size and read the channel data
			sampleRate := int(e.mainlinePacketHeader.sampleRate)
			dataBuf := make([]byte, e.mainlinePacketChannel.length(sampleRate))
			_, err = deps.Transport.Read(dataBuf, time.Second, false)
			if err != nil {
				continue
			}
			err = e.mainlinePacketChannel.decode(dataBuf, sampleRate)
			if err != nil {
				deps.Health.Errors++
				continue
			}

			// Get tail section data, check tail bytes of the packet
			tailBuf := make([]byte, e.mainlinePacketTail.length()+len(mainline_packet_frame_tail))
			_, err = deps.Transport.Read(tailBuf, time.Second, false)
			if err != nil {
				continue
			}
			frameTailSliceIndex := len(tailBuf) - len(mainline_packet_frame_tail)
			if !bytes.Equal(tailBuf[frameTailSliceIndex:], mainline_packet_frame_tail) {
				deps.Health.Errors++
				continue
			}
			err = e.mainlinePacketTail.decode(tailBuf[:frameTailSliceIndex])
			if err != nil {
				deps.Health.Errors++
				continue
			}

			// Get current timestamp
			if e.mainlinePacketHeader.timestamp == 0 {
				t, err := deps.FallbackTime.Get()
				if err != nil {
					continue
				}
				e.mainlinePacketHeader.timestamp = t.UnixMilli()
			}

			// Publish the data to the message bus
			deps.Health.SampleRate = sampleRate
			finalPacket := ExplorerData{
				SampleRate: sampleRate,
				Timestamp:  e.mainlinePacketHeader.timestamp,
				Z_Axis:     e.mainlinePacketChannel.z_axis,
				E_Axis:     e.mainlinePacketChannel.e_axis,
				N_Axis:     e.mainlinePacketChannel.n_axis,
			}
			deps.messageBus.Publish("explorer", &finalPacket)

			deps.Health.UpdatedAt = time.UnixMilli(e.mainlinePacketHeader.timestamp)
			deps.Health.Received++
		}
	}
}

func (e *ExplorerDriverImpl) readerDaemon(deps *ExplorerDependency) {
	if deps.Config.LegacyMode {
		e.handleReadLegacyPacket(deps)
	} else {
		e.handleReadMainlinePacket(deps)
	}
}

func (e *ExplorerDriverImpl) IsAvailable(deps *ExplorerDependency) bool {
	buf := make([]byte, 128)
	_, err := deps.Transport.Read(buf, 2*time.Second, true)
	return err == nil
}

func (e *ExplorerDriverImpl) Init(deps *ExplorerDependency) error {
	currentTime, err := deps.FallbackTime.Get()
	if err != nil {
		return err
	}

	deps.Health.StartTime = currentTime
	deps.subscribers = cmap.New[ExplorerEventHandler]()
	deps.messageBus = messagebus.New(1024)
	deps.Config.DeviceId = math.MaxUint32

	// Get device ID in EEPROM
	if !deps.Config.LegacyMode {
		readTimeout := 5 * time.Second
		startTime := time.Now()
		for time.Since(startTime) < readTimeout {
			ok, _ := deps.Transport.Filter(mainline_packet_frame_header, 2*time.Second)
			if !ok {
				continue
			}
			headerBuf := make([]byte, e.mainlinePacketHeader.length())
			_, err := deps.Transport.Read(headerBuf, time.Second, false)
			if err != nil {
				continue
			}
			err = e.mainlinePacketHeader.decode(headerBuf)
			if err != nil {
				continue
			}
			deps.Config.DeviceId = e.mainlinePacketHeader.deviceId
			break
		}
		if time.Since(startTime) >= readTimeout {
			return errors.New("failed to get device ID, please check the device")
		}
	}

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
	deps.messageBus.Unsubscribe("explorer", fn)
	deps.subscribers.Remove(clientId)
	return nil
}
