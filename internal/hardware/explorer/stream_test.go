package explorer

import (
	"context"
	"encoding/binary"
	"errors"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/anyshake/observer/pkg/timesource"
)

type scriptFrame struct {
	data    []byte
	timeout bool
	err     error
}

type scriptTransport struct {
	packets   chan []byte
	frames    chan scriptFrame
	stop      chan struct{}
	once      sync.Once
	delay     time.Duration
	readErr   error
	done      bool
	immediate bool
}

func newScriptTransport() *scriptTransport {
	return &scriptTransport{
		packets: make(chan []byte, 64),
		frames:  make(chan scriptFrame),
		stop:    make(chan struct{}),
	}
}

func (s *scriptTransport) Open() error  { return nil }
func (s *scriptTransport) Flush() error { return nil }
func (s *scriptTransport) Close() error {
	s.once.Do(func() { close(s.stop) })
	return nil
}
func (s *scriptTransport) Write(buf []byte) (int, error)  { return len(buf), nil }
func (s *scriptTransport) GetLatency(int) time.Duration   { return 0 }
func (s *scriptTransport) SetTimeout(time.Duration) error { return nil }

func (s *scriptTransport) Read(buf []byte) (int, error) {
	if s.immediate {
		select {
		case <-s.stop:
			return 0, errors.New("closed")
		case packet := <-s.packets:
			return copy(buf, packet), nil
		default:
			return 0, nil
		}
	}
	if s.delay > 0 {
		time.Sleep(s.delay)
		select {
		case <-s.stop:
			return 0, errors.New("closed")
		case packet := <-s.packets:
			return copy(buf, packet), nil
		default:
			return 0, nil
		}
	}
	select {
	case <-s.stop:
		return 0, errors.New("closed")
	case packet := <-s.packets:
		if s.readErr != nil {
			return 0, s.readErr
		}
		return copy(buf, packet), nil
	}
}

func (s *scriptTransport) ReadUntil(ctx context.Context, _ int, done func(*[]byte, *time.Time) bool, _ time.Duration) ([]byte, bool, time.Duration, error) {
	if !s.done && done != nil {
		s.done = true
		short := []byte{0x01}
		updated := time.Time{}
		_ = done(&short, &updated)
		tailer := []byte{0x00, 0x11, 0xEF, 0x10}
		_ = done(&tailer, &updated)
		delim := []byte{0xEF, 0x10, 0x01, 0xFE}
		_ = done(&delim, &updated)
	}
	select {
	case <-ctx.Done():
		return nil, false, 0, ctx.Err()
	case <-s.stop:
		return nil, false, 0, errors.New("closed")
	case frame := <-s.frames:
		return frame.data, frame.timeout, 0, frame.err
	}
}

func ntpFixed(now time.Time) uint64 {
	const nano = uint64(time.Second)
	elapsed := uint64(now.Sub(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)))
	seconds := elapsed / nano
	fraction := ((elapsed % nano) << 32) / nano
	return seconds<<32 | fraction
}

func startTestNTP(t *testing.T) string {
	t.Helper()
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1")})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	go func() {
		buffer := make([]byte, 512)
		for {
			n, addr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				return
			}
			if n < 48 {
				continue
			}
			response := make([]byte, 48)
			response[0] = 0x24
			response[1] = 1
			response[2] = 4
			response[3] = 0xEC
			binary.BigEndian.PutUint32(response[8:12], 1)
			now := time.Now()
			binary.BigEndian.PutUint64(response[16:24], ntpFixed(now.Add(-time.Minute)))
			copy(response[24:32], buffer[40:48])
			stamp := ntpFixed(now)
			binary.BigEndian.PutUint64(response[32:40], stamp)
			binary.BigEndian.PutUint64(response[40:48], stamp)
			_, _ = conn.WriteToUDP(response, addr)
		}
	}()
	return "ntp://127.0.0.1:" + strconv.Itoa(conn.LocalAddr().(*net.UDPAddr).Port)
}

func waitUntil(timeout time.Duration, cond func() bool) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return true
		}
		time.Sleep(5 * time.Millisecond)
	}
	return cond()
}

func buildV1WirePacket(sample int32) []byte {
	const channelSize = 5
	packet := make([]byte, 66)
	packet[0], packet[1] = 0xFC, 0x1B
	var checksum [3]byte
	offset := 2
	for axis := 0; axis < 3; axis++ {
		for i := 0; i < channelSize; i++ {
			binary.LittleEndian.PutUint32(packet[offset+i*4:], uint32(sample))
		}
		for _, b := range packet[offset : offset+channelSize*4] {
			checksum[axis] ^= b
		}
		offset += channelSize * 4
	}
	copy(packet[offset:], checksum[:])
	return packet
}

func buildV2WirePacket(timestamp int64, variable uint32, sample int32) []byte {
	const channelSize = 5
	packet := make([]byte, 75)
	packet[0], packet[1] = 0xFA, 0xDE
	binary.LittleEndian.PutUint64(packet[2:10], uint64(timestamp))
	binary.LittleEndian.PutUint32(packet[10:14], variable)
	offset := 14
	for axis := 0; axis < 3; axis++ {
		for i := 0; i < channelSize; i++ {
			binary.LittleEndian.PutUint32(packet[offset:], uint32(sample))
			offset += 4
		}
	}
	var checksum byte
	for _, b := range packet[2 : len(packet)-1] {
		checksum ^= b
	}
	packet[len(packet)-1] = checksum
	return packet
}

func buildV3Frame(timestamp int64, config, variable uint32, channel []byte) []byte {
	body := make([]byte, 16+len(channel)+1)
	binary.LittleEndian.PutUint64(body[0:8], uint64(timestamp))
	binary.LittleEndian.PutUint32(body[8:12], config)
	binary.LittleEndian.PutUint32(body[12:16], variable)
	copy(body[16:], channel)
	var checksum byte
	for _, b := range body[:len(body)-1] {
		checksum ^= b
	}
	body[len(body)-1] = checksum
	return append(body, 0xEF, 0x10, 0x01, 0xFE)
}

func TestV1OpenDecodesWirePackets(t *testing.T) {
	transport := newScriptTransport()
	transport.delay = 100 * time.Millisecond
	protocol := &ExplorerProtoImplV1{
		ChannelCodes:    []string{"HNZ", "HNE", "HNN"},
		ExplorerOptions: ExplorerOptions{Protocol: "v1", Model: "models/E-C111G", Latitude: 25, Longitude: 121, Elevation: 10},
		NtpOptions:      NtpOptions{Pool: []string{startTestNTP(t)}, Retry: 0, ReadTimeout: 1},
		Logger:          testLogger(nil),
		TimeSource:      timesource.New(time.Now),
		Transport:       transport,
	}
	t.Cleanup(func() { _ = protocol.Close() })

	ctx, cancel, err := protocol.Open(context.Background())
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer cancel()
	if ctx.Err() != nil {
		t.Fatalf("Open() returned a canceled context: %v", ctx.Err())
	}

	archived := make(chan Event, 4)
	realtime := make(chan Event, 16)
	if err := protocol.Subscribe("archive", func(event Event) { archived <- event }); err != nil {
		t.Fatal(err)
	}
	if err := protocol.SubscribeRealtime("live", func(event Event) { realtime <- event }); err != nil {
		t.Fatal(err)
	}
	transport.packets <- make([]byte, 66)
	wire := buildV1WirePacket(7)
	for i := 0; i < 10; i++ {
		transport.packets <- wire
	}

	select {
	case event := <-realtime:
		if event.SampleRate != 50 || len(event.ChannelData) != 3 || len(event.ChannelData[0].Data) != 5 || event.ChannelData[0].Data[0] != 7 || event.ChannelData[0].ChannelCode != "HNZ" {
			t.Fatalf("realtime event = %#v", event)
		}
	case <-time.After(3 * time.Second):
		config := protocol.GetConfig()
		t.Fatalf("timed out waiting for realtime data, sample rate %d", config.GetSampleRate())
	}
	select {
	case event := <-archived:
		if len(event.ChannelData) != 3 || len(event.ChannelData[2].Data) != 50 || event.ChannelData[2].ChannelCode != "HNN" {
			t.Fatalf("archived event channels = %d samples %d", len(event.ChannelData), len(event.ChannelData[0].Data))
		}
	case <-time.After(3 * time.Second):
		config := protocol.GetConfig()
		status := protocol.GetStatus()
		t.Fatalf("timed out waiting for archived data, sample rate %d frames %d", config.GetSampleRate(), status.GetFrames())
	}

	corrupt := append([]byte(nil), wire...)
	corrupt[2] ^= 0xFF
	transport.packets <- corrupt
	if !waitUntil(time.Second, func() bool {
		status := protocol.GetStatus()
		return status.GetErrors() > 0
	}) {
		t.Fatal("corrupt packet did not increment errors")
	}
	config := protocol.GetConfig()
	if protocol.GetDeviceId() != "012F81AC" || config.GetProtocol() != "v1" || config.GetModel() != "E-C111G" {
		t.Fatalf("device = %s protocol=%s model=%s", protocol.GetDeviceId(), config.GetProtocol(), config.GetModel())
	}
	transport.delay = 0
	transport.immediate = true
	time.Sleep(200 * time.Millisecond)
}

func TestV1OpenFlushFailure(t *testing.T) {
	transport := &fakeTransport{flushErr: errors.New("flush failed")}
	protocol := &ExplorerProtoImplV1{
		Logger:     testLogger(nil),
		TimeSource: timesource.New(time.Now),
		Transport:  transport,
		NtpOptions: NtpOptions{Pool: []string{startTestNTP(t)}, Retry: 0, ReadTimeout: 1},
	}
	_, _, err := protocol.Open(context.Background())
	if err == nil || !strings.Contains(err.Error(), "flush failed") || transport.closed != 1 || transport.flushed != 1 {
		t.Fatalf("cleanup flushed=%d closed=%d err=%v", transport.flushed, transport.closed, err)
	}
}

func TestV2OpenGNSSStream(t *testing.T) {
	transport := newScriptTransport()
	base := time.Unix(1_700_000_000, 0)
	protocol := &ExplorerProtoImplV2{
		ChannelCodes:    []string{"HNZ", "HNE", "HNN"},
		ExplorerOptions: ExplorerOptions{Protocol: "v2", Model: "E-C111G", Latitude: 25, Longitude: 121, Elevation: 10},
		NtpOptions:      NtpOptions{Pool: []string{startTestNTP(t)}, Retry: 0, ReadTimeout: 1},
		Logger:          testLogger(nil),
		TimeSource:      timesource.New(func() time.Time { return base }),
		Transport:       transport,
	}
	t.Cleanup(func() { _ = protocol.Close() })

	openErr := make(chan error, 1)
	go func() {
		_, _, err := protocol.Open(context.Background())
		openErr <- err
	}()

	const start = int64(1_700_000_000_000)
	deviceWord := uint32(0x80000000) | 0x012F81AC
	for i := 0; i < 10; i++ {
		transport.packets <- buildV2WirePacket(start, deviceWord, 4)
	}
	if !waitUntil(2*time.Second, func() bool { return protocol.GetDeviceId() == "012F81AC" }) {
		t.Fatalf("device id = %s stable=%v", protocol.GetDeviceId(), protocol.isDataStreamStable)
	}
	transport.packets <- buildV2WirePacket(start+1000, math.Float32bits(25.5), 4)
	if !waitUntil(2*time.Second, func() bool {
		_, err := protocol.deviceVariable.GetLatitude(false)
		return err == nil
	}) {
		t.Fatal("latitude was not decoded")
	}
	transport.packets <- buildV2WirePacket(start+2000, math.Float32bits(121.5), 4)
	if !waitUntil(2*time.Second, func() bool {
		_, err := protocol.deviceVariable.GetLongitude(false)
		return err == nil
	}) {
		t.Fatal("longitude was not decoded")
	}
	transport.packets <- buildV2WirePacket(start+3000, math.Float32bits(10.5), 4)
	if !waitUntil(2*time.Second, func() bool { return protocol.variableAllSet }) {
		t.Fatal("variables did not settle")
	}
	transport.packets <- buildV2WirePacket(start+4000, deviceWord, 4)
	if !waitUntil(2*time.Second, func() bool {
		config := protocol.GetConfig()
		return config.GetSampleRate() == 5
	}) {
		config := protocol.GetConfig()
		t.Fatalf("sample rate = %d", config.GetSampleRate())
	}
	transport.packets <- buildV2WirePacket(start+5000, deviceWord, 4)
	select {
	case err := <-openErr:
		if err != nil {
			t.Fatalf("Open() error = %v", err)
		}
	case <-time.After(2 * time.Second):
		config := protocol.GetConfig()
		t.Fatalf("Open() did not become ready, rate=%d set=%v", config.GetSampleRate(), protocol.variableAllSet)
	}

	archived := make(chan Event, 4)
	if err := protocol.Subscribe("archive", func(event Event) { archived <- event }); err != nil {
		t.Fatal(err)
	}
	transport.packets <- buildV2WirePacket(start+7000, deviceWord, 9)
	transport.packets <- buildV2WirePacket(start+8000, deviceWord, 9)
	transport.packets <- buildV2WirePacket(start+9000, deviceWord, 9)
	select {
	case event := <-archived:
		if !event.GNSSAvailable || len(event.ChannelData) != 3 || len(event.ChannelData[0].Data) != 5 || event.ChannelData[0].ChannelCode != "HNZ" {
			t.Fatalf("archived event = %#v", event)
		}
	case <-time.After(2 * time.Second):
		status := protocol.GetStatus()
		t.Fatalf("timed out waiting for archived GNSS data, errors=%d", status.GetErrors())
	}

	transport.packets <- buildV2WirePacket(start+9000, deviceWord, 9)
	transport.packets <- buildV2WirePacket(start+9100, deviceWord, 9)
	bad := buildV2WirePacket(start+9500, deviceWord, 1)
	bad[len(bad)-1] ^= 0xFF
	transport.packets <- bad
	if !waitUntil(time.Second, func() bool {
		status := protocol.GetStatus()
		return status.GetErrors() > 0
	}) {
		t.Fatal("corrupt packet did not increment errors")
	}
}

func TestV3OpenGNSSStream(t *testing.T) {
	transport := newScriptTransport()
	base := time.Unix(1_700_000_000, 0)
	protocol := &ExplorerProtoImplV3{
		ChannelCodes:    []string{"HNZ", "HNE", "HNN"},
		ExplorerOptions: ExplorerOptions{Protocol: "v3", Model: "E-C111G", Latitude: 25, Longitude: 121, Elevation: 10},
		NtpOptions:      NtpOptions{Pool: []string{startTestNTP(t)}, Retry: 0, ReadTimeout: 1},
		Logger:          testLogger(nil),
		TimeSource:      timesource.New(func() time.Time { return base }),
		Transport:       transport,
	}
	t.Cleanup(func() { _ = protocol.Close() })

	openErr := make(chan error, 1)
	go func() {
		_, _, err := protocol.Open(context.Background())
		openErr <- err
	}()
	time.Sleep(20 * time.Millisecond)

	const config = uint32(1<<26) | uint32(1<<24) | 0x1F
	channel := []byte{0x2C, 0x01}
	send := func(timestamp int64, variable uint32) {
		t.Helper()
		transport.frames <- scriptFrame{data: buildV3Frame(timestamp, config, variable, channel)}
	}
	transport.frames <- scriptFrame{timeout: true}
	transport.frames <- scriptFrame{}
	transport.frames <- scriptFrame{data: []byte{0, 0, 0, 0}}
	send(0, 0x012F81AC)
	if !waitUntil(2*time.Second, func() bool { return protocol.GetDeviceId() == "012F81AC" }) {
		t.Fatalf("device id = %s", protocol.GetDeviceId())
	}
	send(1000, math.Float32bits(25.5))
	send(2000, math.Float32bits(121.5))
	send(3000, math.Float32bits(10.5))
	send(4000, math.Float32bits(18.5))
	if !waitUntil(2*time.Second, func() bool { return protocol.variableAllSet }) {
		t.Fatal("variables did not settle")
	}
	if got, err := protocol.GetTemperature(); err != nil || got != float64(float32(18.5)) {
		t.Fatalf("temperature = %v, %v", got, err)
	}
	send(5000, 0)
	select {
	case err := <-openErr:
		if err != nil {
			t.Fatalf("Open() error = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Open() did not become ready")
	}

	archived := make(chan Event, 2)
	realtime := make(chan Event, 16)
	if err := protocol.Subscribe("archive", func(event Event) { archived <- event }); err != nil {
		t.Fatal(err)
	}
	if err := protocol.SubscribeRealtime("live", func(event Event) { realtime <- event }); err != nil {
		t.Fatal(err)
	}
	for timestamp := int64(6000); timestamp < 16000; timestamp += 1000 {
		send(timestamp, 0)
	}
	select {
	case event := <-realtime:
		if len(event.ChannelData) != 1 || event.ChannelData[0].ChannelCode != "HNZ" || event.ChannelData[0].Data[0] != 300 {
			t.Fatalf("realtime event = %#v", event)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for realtime data")
	}
	select {
	case event := <-archived:
		if event.SampleRate != 10 || len(event.ChannelData[0].Data) != 10 {
			t.Fatalf("archived event = %#v", event)
		}
	case <-time.After(2 * time.Second):
		config := protocol.GetConfig()
		t.Fatalf("timed out waiting for archived data, samples=%d rate=%d", protocol.collectedSamples, config.GetSampleRate())
	}

	send(17000, 0)
	send(17100, 0)
	bad := buildV3Frame(18000, config, 0, channel)
	bad[0] ^= 0xFF
	transport.frames <- scriptFrame{data: bad}
	if !waitUntil(time.Second, func() bool {
		status := protocol.GetStatus()
		return status.GetErrors() > 0
	}) {
		t.Fatal("corrupt packet did not increment errors")
	}
}

func TestV3OpenReadErrorCancels(t *testing.T) {
	transport := newScriptTransport()
	protocol := &ExplorerProtoImplV3{
		Logger:     testLogger(nil),
		TimeSource: timesource.New(time.Now),
		Transport:  transport,
		NtpOptions: NtpOptions{Pool: []string{startTestNTP(t)}, Retry: 0, ReadTimeout: 1},
	}
	t.Cleanup(func() { _ = protocol.Close() })
	openErr := make(chan error, 1)
	go func() {
		_, _, err := protocol.Open(context.Background())
		openErr <- err
	}()
	transport.frames <- scriptFrame{err: errors.New("read failed")}
	select {
	case err := <-openErr:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Open() error = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("read error did not cancel startup")
	}
}

func TestV2OpenNTPSync(t *testing.T) {
	transport := newScriptTransport()
	base := time.Unix(1_700_000_000, 0)
	protocol := &ExplorerProtoImplV2{
		ChannelCodes:    []string{"HNZ", "HNE", "HNN"},
		ExplorerOptions: ExplorerOptions{Protocol: "v2", Model: "E-C111G", Latitude: 25.5, Longitude: 121.5, Elevation: 10.5},
		NtpOptions:      NtpOptions{Pool: []string{startTestNTP(t)}, Retry: 0, ReadTimeout: 1},
		Logger:          testLogger(nil),
		TimeSource:      timesource.New(func() time.Time { return base }),
		Transport:       transport,
	}
	t.Cleanup(func() { _ = protocol.Close() })
	openErr := make(chan error, 1)
	go func() {
		_, _, err := protocol.Open(context.Background())
		openErr <- err
	}()

	const shifted = int64(1_699_999_900_000)
	for i := 1; i <= 9; i++ {
		transport.packets <- buildV2WirePacket(shifted+int64(i)*100, 0x012F81AC, 3)
	}
	transport.packets <- buildV2WirePacket(shifted, 0x012F81AC, 3)
	for i := 0; i < 10; i++ {
		transport.packets <- buildV2WirePacket(shifted, 0x012F81AC, 3)
	}
	if !waitUntil(2*time.Second, func() bool { return protocol.GetDeviceId() == "012F81AC" }) {
		t.Fatalf("device id %s stable=%v", protocol.GetDeviceId(), protocol.isDataStreamStable)
	}
	transport.packets <- buildV2WirePacket(shifted+1000, 0, 3)
	if !waitUntil(2*time.Second, func() bool {
		config := protocol.GetConfig()
		return config.GetSampleRate() > 0
	}) {
		config := protocol.GetConfig()
		t.Fatalf("sample rate %d diff %d", config.GetSampleRate(), protocol.timeDiff4NonGnssMode)
	}
	transport.packets <- buildV2WirePacket(shifted+2000, 0, 3)
	select {
	case err := <-openErr:
		if err != nil {
			t.Fatalf("Open() error = %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("NTP synchronization did not finish")
	}
	transport.packets <- buildV2WirePacket(shifted+2100, 0, 3)
	transport.packets <- buildV2WirePacket(shifted+9100, 0, 3)
	if !waitUntil(time.Second, func() bool { return !protocol.variableAllSet }) {
		t.Fatal("time jump did not reset variables")
	}
}

func TestV3OpenNTPSync(t *testing.T) {
	transport := newScriptTransport()
	base := time.Unix(1_700_000_000, 0)
	protocol := &ExplorerProtoImplV3{
		ChannelCodes:    []string{"HNZ"},
		ExplorerOptions: ExplorerOptions{Protocol: "v3", Model: "E-C111G", Latitude: 25.5, Longitude: 121.5, Elevation: 10.5},
		NtpOptions:      NtpOptions{Pool: []string{startTestNTP(t)}, Retry: 0, ReadTimeout: 1},
		Logger:          testLogger(nil),
		TimeSource:      timesource.New(func() time.Time { return base }),
		Transport:       transport,
	}
	t.Cleanup(func() { _ = protocol.Close() })
	openErr := make(chan error, 1)
	go func() {
		_, _, err := protocol.Open(context.Background())
		openErr <- err
	}()

	const config = uint32(1<<24) | 0x1F
	channel := []byte{0x01, 0x00}
	const start = int64(1_600_000_000_000)
	send := func(timestamp int64, variable uint32) {
		t.Helper()
		transport.frames <- scriptFrame{data: buildV3Frame(timestamp, config, variable, channel)}
	}
	for i := 0; i < 10; i++ {
		send(start, 0x012F81AC)
	}
	if !waitUntil(2*time.Second, func() bool { return protocol.GetDeviceId() == "012F81AC" }) {
		t.Fatalf("device id %s stable=%v", protocol.GetDeviceId(), protocol.isDataStreamStable)
	}
	send(start+1000, math.Float32bits(1))
	send(start+2000, math.Float32bits(2))
	send(start+3000, math.Float32bits(3))
	send(start+4000, math.Float32bits(4))
	if !waitUntil(2*time.Second, func() bool { return protocol.variableAllSet }) {
		t.Fatal("variables did not settle")
	}
	send(start+5000, 0)
	select {
	case err := <-openErr:
		if err != nil {
			t.Fatalf("Open() error = %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("NTP synchronization did not finish")
	}
	send(start+20000, 0)
	if !waitUntil(time.Second, func() bool { return !protocol.variableAllSet }) {
		t.Fatal("time jump did not reset variables")
	}
}
