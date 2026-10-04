package explorer

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestV3DeviceConfig(t *testing.T) {
	t.Parallel()

	protocol := &ExplorerProtoImplV3{}
	intervals := []time.Duration{100, 200, 500, 1000}
	rates := []int{10, 20, 50, 100, 200, 250, 500, 1000}
	for intervalCode, interval := range intervals {
		for rateCode, rate := range rates {
			t.Run(fmt.Sprintf("%dms_%dHz", interval, rate), func(t *testing.T) {
				// Enable all channels and variables to catch overlapping bit masks.
				config := uint32(intervalCode)<<30 | uint32(rateCode)<<27 | 0x03ffffff
				if got := protocol.parsePacketInterval(config); got != interval*time.Millisecond {
					t.Errorf("packet interval = %v, want %v", got, interval*time.Millisecond)
				}
				if got := protocol.parseSampleRate(config); got != rate {
					t.Errorf("sample rate = %d, want %d", got, rate)
				}
				if protocol.parseGnssAvailibility(config) {
					t.Error("GNSS enabled with its bit cleared")
				}
				if !protocol.parseGnssAvailibility(config | 1<<26) {
					t.Error("GNSS disabled with its bit set")
				}
				chunkLength, totalSize, channels := protocol.getChannelSize(config)
				wantSamples := int(interval) * rate / 1000
				if chunkLength != wantSamples || totalSize != wantSamples*8*4 || len(channels) != 8 {
					t.Fatalf("channel layout = (%d, %d, %d), want (%d, %d, 8)", chunkLength, totalSize, len(channels), wantSamples, wantSamples*8*4)
				}
			})
		}
	}
}

func TestV3MixedChannelLayout(t *testing.T) {
	t.Parallel()

	protocol := &ExplorerProtoImplV3{ChannelCodes: []string{"HNZ", "HNE", "HNN"}}
	// 200 ms at 100 Hz, with channels 1, 3 and 8 enabled.
	config := uint32(1<<30 | 3<<27 | 1<<24 | 2<<20 | 3<<10)
	chunkLength, totalSize, channels := protocol.getChannelSize(config)
	if chunkLength != 20 || totalSize != 180 {
		t.Fatalf("channel size = (%d, %d), want (20, 180)", chunkLength, totalSize)
	}
	want := []*ChannelData{
		{ChannelCode: "HNZ", ChannelId: 1, ByteSize: 2, DataType: "int16"},
		{ChannelCode: "HNN", ChannelId: 3, ByteSize: 3, DataType: "int24"},
		{ChannelCode: "CH8", ChannelId: 8, ByteSize: 4, DataType: "int32"},
	}
	if !reflect.DeepEqual(channels, want) {
		t.Fatalf("channels = %#v, want %#v", channels, want)
	}
	_, totalSize, channels = protocol.getChannelSize(0)
	if totalSize != 0 || len(channels) != 0 {
		t.Fatalf("disabled channels = (%d, %#v), want no data", totalSize, channels)
	}
}

func TestV3ChannelDataSignedSamples(t *testing.T) {
	t.Parallel()

	channels := []*ChannelData{
		{DataType: "int16"},
		{DataType: "int24"},
		{DataType: "int32"},
	}
	// Each channel contains zero, -1, minimum and maximum, in wire byte order.
	packet := []byte{
		0x00, 0x00, 0xff, 0xff, 0x00, 0x80, 0xff, 0x7f,
		0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0x00, 0x00, 0x80, 0xff, 0xff, 0x7f,
		0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,
		0x00, 0x00, 0x00, 0x80, 0xff, 0xff, 0xff, 0x7f,
	}
	want := [][]int32{
		{0, -1, -32768, 32767},
		{0, -1, -8388608, 8388607},
		{0, -1, -2147483648, 2147483647},
	}
	protocol := &ExplorerProtoImplV3{}
	protocol.getChannelData(channels, packet, 4)
	for i, channel := range channels {
		if !reflect.DeepEqual(channel.Data, want[i]) {
			t.Errorf("%s samples = %v, want %v", channel.DataType, channel.Data, want[i])
		}
	}
	protocol.getChannelData(channels, make([]byte, 9), 1)
	for _, channel := range channels {
		if !reflect.DeepEqual(channel.Data, []int32{0}) {
			t.Errorf("%s retained samples from previous packet: %v", channel.DataType, channel.Data)
		}
	}
}

func TestProtocolChecksums(t *testing.T) {
	t.Parallel()

	v2 := &ExplorerProtoImplV2{}
	v3 := &ExplorerProtoImplV3{}
	header := []byte{0xaa, 0x55}
	tests := []struct {
		name    string
		verify  func([]byte) error
		packet  []byte
		wantErr bool
	}{
		{"v2 valid", func(p []byte) error { return v2.verifyChecksum(p, header) }, []byte{0xaa, 0x55, 0x12, 0x34, 0x26}, false},
		{"v2 corrupt", func(p []byte) error { return v2.verifyChecksum(p, header) }, []byte{0xaa, 0x55, 0x13, 0x34, 0x26}, true},
		{"v2 empty", func(p []byte) error { return v2.verifyChecksum(p, header) }, nil, true},
		{"v2 header only", func(p []byte) error { return v2.verifyChecksum(p, header) }, header, true},
		{"v2 shorter than header", func(p []byte) error { return v2.verifyChecksum(p, header) }, []byte{0xaa}, true},
		{"v3 valid", v3.verifyChecksum, []byte{0x12, 0x34, 0x26}, false},
		{"v3 corrupt", v3.verifyChecksum, []byte{0x13, 0x34, 0x26}, true},
		{"v3 empty", v3.verifyChecksum, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.verify(tt.packet); (err != nil) != tt.wantErr {
				t.Fatalf("verifyChecksum() error = %v, want error = %v", err, tt.wantErr)
			}
		})
	}
}

func TestV2TimestampCorrection(t *testing.T) {
	t.Parallel()

	protocol := &ExplorerProtoImplV2{timeDiff4NonGnssMode: 1500}
	if got := protocol.getTimestamp(10000); got != 11500 {
		t.Fatalf("non-GNSS timestamp = %d, want 11500", got)
	}
	protocol.deviceConfig.SetGnssAvailability(true)
	if got := protocol.getTimestamp(10000); got != 10000 {
		t.Fatalf("GNSS timestamp = %d, want 10000", got)
	}
}
