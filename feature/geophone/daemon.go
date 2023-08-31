package geophone

import (
	"fmt"
	"os"
	"time"

	"com.geophone.observer/common/serial"
	"com.geophone.observer/feature"
	t "com.geophone.observer/utils/time"
)

func (g *Geophone) Start(options *feature.FeatureOptions) {
	var (
		device = options.Config.Serial.Device
		baud   = options.Config.Serial.Baud
		length = options.Config.Serial.Length
	)

	port, err := serial.Open(device, baud)
	if err != nil {
		options.OnError(MODULE, options, err)
		os.Exit(1)
	}
	defer serial.Close(port)

	var packet Packet
	options.OnStart(MODULE, options, "service has started")

	lastRead := time.Now().UTC()
	for {
		err := g.Read(port, &packet, length)
		if err != nil {
			serial.Close(port)
			options.OnError(MODULE, options, err)

			port, err = serial.Open(device, baud)
			if err != nil {
				options.OnError(MODULE, options, err)
				os.Exit(1)
			}

			err = g.Reset(port)
			if err != nil {
				options.OnError(MODULE, options, err)
			}

			lastRead = time.Now().UTC()
			continue
		} else {
			options.OnReady(MODULE, options, packet)
		}

		// Reset device if reached TIMEOUT_THRESHOLD
		if t.Diff(time.Now().UTC(), lastRead) >= TIMEOUT_THRESHOLD {
			err := fmt.Errorf("reset due to unusual gap")
			options.OnError(MODULE, options, err)

			err = g.Reset(port)
			if err != nil {
				options.OnError(MODULE, options, err)
			}
		}

		lastRead = time.Now().UTC()
	}
}
