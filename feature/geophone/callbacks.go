package geophone

import (
	"time"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/utils/duration"
	"github.com/anyshake/observer/utils/logger"
	"github.com/anyshake/observer/utils/text"
	"github.com/fatih/color"
)

func (g *Geophone) OnStart(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, text.Concat(v...), color.FgMagenta, false)
}

func (g *Geophone) OnStop(options *feature.FeatureOptions, v ...any) {
	logger.Print(MODULE, text.Concat(v...), color.FgBlue, true)
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
			// Get compensation filter coefficients
			// Set packet timestamp
			options.Status.System.Messages++
			options.Status.Buffer.TS = currentTime.UnixMilli()
			// Apply compensation for EHZ, EHE, EHN channels
			if len(packet.EHZ) > 1 && len(packet.EHE) > 1 && len(packet.EHN) > 1 {
				packet.EHZ = g.applyFilter(packet.EHZ, g.getFilter(v[1].(config.Compensation), len(packet.EHZ)))
				packet.EHE = g.applyFilter(packet.EHE, g.getFilter(v[2].(config.Compensation), len(packet.EHE)))
				packet.EHN = g.applyFilter(packet.EHN, g.getFilter(v[3].(config.Compensation), len(packet.EHN)))
			}
			// Copy buffer and reset
			options.Status.Geophone = *options.Status.Buffer
			options.Status.Buffer.EHZ = []int32{}
			options.Status.Buffer.EHE = []int32{}
			options.Status.Buffer.EHN = []int32{}
			logger.Print(MODULE, "1 full packet received", color.FgGreen, false)
		}
	} else {
		logger.Print(MODULE, "waiting for time alignment", color.FgYellow, false)
	}
}

func (g *Geophone) OnError(options *feature.FeatureOptions, err error) {
	options.Status.System.Errors++
	logger.Print(MODULE, err, color.FgRed, false)
}
