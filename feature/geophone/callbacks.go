package geophone

import (
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
	if !options.Status.ReadyTime.IsZero() {
		var (
			packet         = v[0].(Packet)
			currentTime, _ = duration.Timestamp(options.Status.System.Offset)
		)

		// Appending packet data to buffer
		for i := 0; i < options.Config.Serial.Packet; i++ {
			options.Status.Buffer.EHZ = append(options.Status.Buffer.EHZ, packet.EHZ[i])
			options.Status.Buffer.EHE = append(options.Status.Buffer.EHE, packet.EHE[i])
			options.Status.Buffer.EHN = append(options.Status.Buffer.EHN, packet.EHN[i])
		}

		// Archive approximately 1 second has passed
		timeDiff := duration.Difference(currentTime, options.Status.LastRecvTime)
		if timeDiff >= READY_THRESHOLD {
			// Set packet timestamp, note that the timestamp in buffer is the start of the packet
			options.Status.Buffer.TS = currentTime.UnixMilli() - timeDiff.Milliseconds()
			// Set last received time is the current timestamp
			options.Status.LastRecvTime = currentTime
			options.Status.System.Messages++
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
