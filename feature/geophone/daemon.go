package geophone

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bclswl0827/observer/driver/serial"
	"github.com/bclswl0827/observer/feature"
	"github.com/bclswl0827/observer/utils/duration"
)

func (g *Geophone) Run(options *feature.FeatureOptions, waitGroup *sync.WaitGroup) {
	var (
		device    = options.Config.Serial.Device
		baud      = options.Config.Serial.Baud
		packetLen = options.Config.Serial.Packet
	)

	// Increase wait group counter
	waitGroup.Add(1)
	defer waitGroup.Done()

	// Open serial port
	port, err := serial.Open(device, baud)
	if err != nil {
		g.OnError(options, err)
		os.Exit(1)
	}
	defer serial.Close(port)

	go func() {
		// Initialize geophone packet
		var packet Packet
		g.OnStart(options, "service has started")

		// FIXME: Compensation is in beta
		if options.Config.Geophone.EHZ.Compensation ||
			options.Config.Geophone.EHE.Compensation ||
			options.Config.Geophone.EHN.Compensation {
			g.OnStart(options, "compensation is in beta")
		}

		lastRead := time.Now().UTC()
		for {
			// Read from serial port by channel packet length
			err := g.Read(port, options.Config, &packet, packetLen)
			if err != nil {
				serial.Close(port)
				g.OnError(options, err)
				time.Sleep(time.Millisecond * 100)

				// Reopen serial port
				port, err = serial.Open(device, baud)
				if err != nil {
					g.OnError(options, err)
					os.Exit(1)
				}

				// Reset device after reopen
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
	}()

	// Receive interrupt signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signals
	<-sigCh
	g.OnStop(options, "closing serial connection")
	serial.Close(port)
}
