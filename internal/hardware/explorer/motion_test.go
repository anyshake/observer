package explorer

import (
	"testing"
	"time"

	"github.com/anyshake/observer/pkg/groundmotion"
	"github.com/anyshake/observer/pkg/metadata"
)

func TestPressureChannelMetadata(t *testing.T) {
	g := &groundMotion{
		sensitivities: []*metadata.Sensitivity{{Gain: 1e6, Unit: metadata.UNIT_PRESSURE}},
		channels:      make(map[int]*groundmotion.Channel),
	}
	channels := []ChannelData{
		{ChannelId: 1, ChannelCode: "HDF", Data: []int32{125000, -250000}},
		{ChannelId: 2, ChannelCode: "CH2", Data: []int32{42}},
	}
	g.annotate(time.Now(), 100, channels)
	pressure := channels[0]
	if pressure.InputUnit != "Pa" || pressure.Sensitivity != 1e6 || pressure.Data[0] != 125000 {
		t.Fatalf("pressure channel = %+v", pressure)
	}
	if pressure.PGA != nil || pressure.PGV != nil || pressure.PGD != nil {
		t.Fatal("pressure channel has ground motion measurements")
	}
	if channels[1].InputUnit != "" || channels[1].Sensitivity != 0 {
		t.Fatal("raw channel acquired a physical unit")
	}
}
