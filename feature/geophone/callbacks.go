package geophone

import (
	"time"

	"github.com/bclswl0827/observer/feature"
	"github.com/bclswl0827/observer/utils/duration"
	"github.com/bclswl0827/observer/utils/logger"
	"github.com/bclswl0827/observer/utils/text"
	"github.com/fatih/color"
)

func (g *Geophone) OnStart(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, text.Concat(v...), color.FgMagenta)
}

func (g *Geophone) OnStop(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, text.Concat(v...), color.FgBlue)
}

func (g *Geophone) OnReady(options *feature.FeatureOptions, v ...any) {
	if options.Status.IsReady {
		var (
			packet         = v[0].(Packet)
			lastTime       = time.UnixMilli(options.Status.Buffer.TS).UTC()
			currentTime, _ = duration.Timestamp(options.Status.System.Offset)
		)

		// Appending packet data to buffer
		for i := 0; i < options.Config.Serial.Packet; i++ {
			options.Status.Buffer.EHZ = append(options.Status.Buffer.EHZ, packet.EHZ[i])
			options.Status.Buffer.EHE = append(options.Status.Buffer.EHE, packet.EHE[i])
			options.Status.Buffer.EHN = append(options.Status.Buffer.EHN, packet.EHN[i])
		}

		// Archive approximately 1 second has passed
		timeDiff := duration.Difference(currentTime, lastTime)
		if timeDiff >= READY_THRESHOLD {
			// Set packet timestamp
			options.Status.System.Messages++
			options.Status.Buffer.TS = currentTime.UnixMilli()
			// Apply compensation for EHZ
			if options.Config.Geophone.EHZ.Compensation {
				packet.EHZ = g.Filter(packet.EHZ, &Filter{
					a1: 1.99823115,
					a2: -0.99822469,
					b0: 1.03380975,
					b1: -1.99662644,
					b2: 0.96601161,
				})
			}
			// Apply compensation for EHE
			if options.Config.Geophone.EHE.Compensation {
				packet.EHE = g.Filter(packet.EHE, &Filter{
					a1: 1.99823115,
					a2: -0.99822469,
					b0: 1.03380975,
					b1: -1.99662644,
					b2: 0.96601161,
				})
			}
			// Apply compensation for EHN
			if options.Config.Geophone.EHN.Compensation {
				packet.EHN = g.Filter(packet.EHN, &Filter{
					a1: 1.99823115,
					a2: -0.99822469,
					b0: 1.03380975,
					b1: -1.99662644,
					b2: 0.96601161,
				})
			}
			// Copy buffer and reset
			options.Status.Geophone = *options.Status.Buffer
			options.Status.Buffer.EHZ = []int32{}
			options.Status.Buffer.EHE = []int32{}
			options.Status.Buffer.EHN = []int32{}
			logger.Print(MODULE, "1 full packet received", color.FgGreen)
		}
	} else {
		logger.Print(MODULE, "waiting for time alignment", color.FgYellow)
	}
}

func (g *Geophone) OnError(options *feature.FeatureOptions, err error) {
	options.Status.System.Errors++
	logger.Print(MODULE, err, color.FgRed)
}
