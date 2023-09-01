package geophone

import (
	"fmt"
	"os"
	"time"

	"github.com/bclswl0827/observer/driver/serial"
	"github.com/bclswl0827/observer/feature"
	"github.com/bclswl0827/observer/utils/duration"
)

func (g *Geophone) Start(options *feature.FeatureOptions) {
	var (
		device = options.Config.Serial.Device
		baud   = options.Config.Serial.Baud
		length = options.Config.Serial.Length
	)

	port, err := serial.Open(device, baud)
	if err != nil {
		g.OnError(options, err)
		os.Exit(1)
	}
	defer serial.Close(port)

	var packet Packet
	g.OnStart(options, "service has started")

	lastRead := time.Now().UTC()
	for {
		err := g.Read(port, &packet, length)
		if err != nil {
			serial.Close(port)
			g.OnError(options, err)

			port, err = serial.Open(device, baud)
			if err != nil {
				g.OnError(options, err)
				os.Exit(1)
			}

			err = g.Reset(port)
			if err != nil {
				g.OnError(options, err)
			}

			lastRead = time.Now().UTC()
			continue
		} else {
			g.OnReady(options, packet)
		}

		// Reset device if reached TIMEOUT_THRESHOLD
		if duration.Difference(time.Now().UTC(), lastRead) >= TIMEOUT_THRESHOLD {
			err := fmt.Errorf("reset due to unusual gap")
			g.OnError(options, err)

			err = g.Reset(port)
			if err != nil {
				g.OnError(options, err)
			}
		}

		lastRead = time.Now().UTC()
	}
}
