package explorer

import (
	"reflect"
	"testing"
	"time"
)

func TestNewEventCopiesChannelData(t *testing.T) {
	t.Parallel()

	timestamp := time.Unix(1700000000, 0)
	config := &DeviceConfig{}
	config.SetSampleRate(100)
	config.SetGnssAvailability(true)
	channels := []ChannelData{
		{ChannelCode: "HNZ", ChannelId: 1, ByteSize: 3, DataType: "int24", Data: []int32{1, -2, 3}},
		{ChannelCode: "HNE", ChannelId: 2, ByteSize: 4, DataType: "int32", Data: []int32{4, 5}},
	}
	event := NewEvent(timestamp, config, channels)
	if !event.Timestamp.Equal(timestamp) || event.SampleRate != 100 || !event.GNSSAvailable {
		t.Fatalf("unexpected event metadata: %#v", event)
	}
	if !reflect.DeepEqual(event.ChannelData, channels) {
		t.Fatalf("event channels = %#v, want %#v", event.ChannelData, channels)
	}

	channels[0].ChannelCode = "changed"
	channels[0].Data[0] = 99
	channels[1].Data[0] = 88
	config.SetSampleRate(50)
	config.SetGnssAvailability(false)
	if event.ChannelData[0].ChannelCode != "HNZ" || event.ChannelData[0].Data[0] != 1 || event.ChannelData[1].Data[0] != 4 {
		t.Fatalf("event changed when source buffers were reused: %#v", event.ChannelData)
	}
	if event.SampleRate != 100 || !event.GNSSAvailable {
		t.Fatal("event metadata changed with device config")
	}
	event.ChannelData[0].Data[1] = 77
	if channels[0].Data[1] != -2 {
		t.Fatal("event mutation changed source samples")
	}
}
