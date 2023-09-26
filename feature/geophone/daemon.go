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
		ehz       = options.Config.Geophone.EHZ
		ehe       = options.Config.Geophone.EHE
		ehn       = options.Config.Geophone.EHN
		device    = options.Config.Serial.Device
		baud      = options.Config.Serial.Baud
		packetLen = options.Config.Serial.Packet
	)

	port, err := serial.Open(device, baud)
	if err != nil {
		g.OnError(options, err)
		os.Exit(1)
	}
	defer serial.Close(port)

	var (
		ehzFilter = g.getFilter(ehz.Sensitivity, ehz.Frequency, ehz.Damping)
		eheFilter = g.getFilter(ehe.Sensitivity, ehe.Frequency, ehe.Damping)
		ehnFilter = g.getFilter(ehn.Sensitivity, ehn.Frequency, ehn.Damping)
	)

	var packet Packet
	g.OnStart(options, "service has started")

	lastRead := time.Now().UTC()
	for {
		err := g.Read(port, &packet, packetLen, ehzFilter, eheFilter, ehnFilter)
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
