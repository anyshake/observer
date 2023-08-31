package callbacks

import (
	"time"

	"com.geophone.observer/feature"
	"com.geophone.observer/feature/geophone"
	"com.geophone.observer/utils/logger"
	t "com.geophone.observer/utils/time"
	"github.com/fatih/color"
)

func OnReady(module string, options *feature.FeatureOptions, v ...any) {
	switch module {
	case "archiver":
		logger.Print(module, "1 message has been archived", color.FgGreen)
	case "miniseed":
		logger.Print(module, "1 record has been written", color.FgGreen)
	case "ntpclient":
		options.Status.System.Offset = v[0].(float64)
		options.Status.IsReady = true
		logger.Print(module, "time alignment succeed", color.FgGreen)
	case "geophone":
		if options.Status.IsReady {
			var (
				packet      = v[0].(geophone.Packet)
				currentTime = time.Now().UTC()
				lastTime    = time.UnixMilli(options.Status.Buffer.TS).UTC()
			)

			// Appending packet data to buffer
			for i := 0; i < options.Config.Serial.Length; i++ {
				options.Status.Buffer.EHZ = append(options.Status.Buffer.EHZ, packet.EHZ[i])
				options.Status.Buffer.EHE = append(options.Status.Buffer.EHE, packet.EHE[i])
				options.Status.Buffer.EHN = append(options.Status.Buffer.EHN, packet.EHN[i])
			}

			// Archive approximately 1 second has passed
			timeDiff := t.Diff(currentTime, lastTime)
			if timeDiff >= geophone.READY_THRESHOLD {
				// Set packet timestamp
				options.Status.System.Messages++
				options.Status.Buffer.TS = currentTime.UnixMilli()
				logger.Print(module, "1 full packet received", color.FgGreen)
				// Copy buffer and reset
				options.Status.Geophone = *options.Status.Buffer
				options.Status.Buffer.EHZ = []int32{}
				options.Status.Buffer.EHE = []int32{}
				options.Status.Buffer.EHN = []int32{}
			}
		} else {
			logger.Print(module, "waiting for time alignment", color.FgYellow)
		}
	}
}
