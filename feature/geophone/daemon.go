package geophone

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/anyshake/observer/driver/serial"
	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/utils/duration"
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

	g.Ticker = time.NewTicker(READY_THRESHOLD)
	defer g.Ticker.Stop()

	go func() {
		// Initialize geophone packet
		var packet Packet
		g.OnStart(options, "service has started")

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

	go func() {
		for {
			<-g.Ticker.C
			currentTime, _ := duration.Timestamp(options.Status.System.Offset)
			timeDiff := duration.Difference(currentTime, options.Status.LastRecvTime)
			// Set packet timestamp, note that the timestamp in buffer is the start of the packet
			options.Status.Buffer.TS = currentTime.UnixMilli() - timeDiff.Milliseconds()
			// Set last received time is the current timestamp
			options.Status.LastRecvTime = currentTime
			options.Status.System.Messages++
			// Copy and reset buffer
			options.Status.Geophone = *options.Status.Buffer
			options.Status.Buffer.EHZ = []int32{}
			options.Status.Buffer.EHE = []int32{}
			options.Status.Buffer.EHN = []int32{}
			g.OnReady(options)
		}
	}()

	// Receive interrupt signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signals
	<-sigCh
	g.OnStop(options, "closing serial connection")
}
