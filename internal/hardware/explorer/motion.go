package explorer

import (
	"sync"
	"time"

	"github.com/anyshake/observer/pkg/groundmotion"
	"github.com/anyshake/observer/pkg/metadata"
)

func cloneFloat(value *float64) *float64 {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

// groundMotion computes per-channel PGA/PGV/PGD for the normal (1 s) stream.
// Sensitivities come from the model metadata, indexed by 1-based channel ID.
type groundMotion struct {
	mu            sync.Mutex
	sensitivities []*metadata.Sensitivity
	channels      map[int]*groundmotion.Channel
}

func newGroundMotion(model string) *groundMotion {
	sensitivities, _ := metadata.Sensitivities(model)
	return &groundMotion{
		sensitivities: sensitivities,
		channels:      make(map[int]*groundmotion.Channel),
	}
}

// annotate fills the motion fields of channelData in place. Channels without
// metadata, or whose filters are still settling, keep nil values.
func (g *groundMotion) annotate(start time.Time, sampleRate int, channelData []ChannelData) {
	if g == nil {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for i := range channelData {
		ch := &channelData[i]
		ch.PGA, ch.PGV, ch.PGD = nil, nil, nil
		calculator, ok := g.channels[ch.ChannelId]
		if !ok {
			if idx := ch.ChannelId - 1; idx >= 0 && idx < len(g.sensitivities) && g.sensitivities[idx] != nil {
				sensitivity := g.sensitivities[idx]
				calculator = groundmotion.NewChannel(sensitivity.Gain, sensitivity.Unit).
					WithResponse(sensitivity.Poles, sensitivity.Frequency)
			}
			g.channels[ch.ChannelId] = calculator
		}
		if calculator == nil {
			continue
		}
		peak, ready := calculator.Process(start, sampleRate, ch.Data)
		if !ready {
			continue
		}
		ch.PGA, ch.PGV, ch.PGD = &peak.PGA, &peak.PGV, &peak.PGD
	}
}
