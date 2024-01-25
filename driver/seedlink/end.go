package seedlink

import (
	"net"
	"time"

	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/publisher"
	"github.com/anyshake/observer/utils/text"
)

type END struct{}

// Callback of "END" command, implements SeedLinkCommandCallback interface
func (*END) Callback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, streamer SeedLinkStreamer, conn net.Conn, args ...string) error {
	cl.StreamMode = true // Enter stream mode
	var (
		seqNum    int64 = 0
		channels        = cl.Channels
		location        = cl.Location
		endTime         = cl.EndTime
		startTime       = cl.StartTime
		station         = text.TruncateString(cl.Station, 5)
		network         = text.TruncateString(cl.Network, 2)
	)

	if startTime.IsZero() {
		_, err := conn.Write([]byte(RES_ERR))
		return err
	}

	// Send data in buffer
	for _, buffer := range sl.SeedLinkBuffer.Data {
		chMap := map[string]publisher.Int32Array{
			"EHZ": buffer.EHZ, "EHE": buffer.EHE, "EHN": buffer.EHN,
		}
		for _, channel := range channels {
			if data, ok := chMap[channel]; ok {
				bufTime := time.UnixMilli(buffer.TS)
				if bufTime.After(startTime) && bufTime.Before(endTime) {
					dataBytes, err := CreateSLPacket(data, buffer.TS, seqNum, network, station, channel, location)
					if err != nil {
						return err
					}

					if len(dataBytes) > 0 {
						_, err = conn.Write(dataBytes)
						if err != nil {
							return err
						}

						seqNum++
					}
				}
			}
		}
	}

	// Subscribe to the publisher
	go publisher.Subscribe(
		&options.Status.Geophone, &cl.StreamMode,
		func(gp *publisher.Geophone) error {
			return streamer(gp, conn, channels, network, station, location, &seqNum)
		},
	)

	return nil
}

// Fallback of "END" command, implements SeedLinkCommandCallback interface
func (*END) Fallback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, conn net.Conn, args ...string) {
	conn.Close()
}
