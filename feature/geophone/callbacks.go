package geophone

import (
	"github.com/anyshake/observer/feature"
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
	if len(v) == 0 {
		logger.Print(MODULE, "1 full packet received", color.FgGreen, false)
		return
	}
	if !options.Status.ReadyTime.IsZero() {
		// Appending packet data to buffer
		packet := v[0].(Packet)
		for i := 0; i < options.Config.Serial.Packet; i++ {
			options.Status.Buffer.EHZ = append(options.Status.Buffer.EHZ, packet.EHZ[i])
			options.Status.Buffer.EHE = append(options.Status.Buffer.EHE, packet.EHE[i])
			options.Status.Buffer.EHN = append(options.Status.Buffer.EHN, packet.EHN[i])
		}
	} else {
		logger.Print(MODULE, "waiting for time alignment", color.FgYellow, false)
	}
}

func (g *Geophone) OnError(options *feature.FeatureOptions, err error) {
	options.Status.System.Errors++
	logger.Print(MODULE, err, color.FgRed, false)
}
