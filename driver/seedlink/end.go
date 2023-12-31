package seedlink

import (
	"net"

	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/publisher"
	"github.com/anyshake/observer/utils/text"
)

type END struct{}

// Callback of "END" command, implements SeedLinkCommandCallback interface
func (*END) Callback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, streamer SeedLinkStreamer, conn net.Conn, args ...string) error {
	cl.StreamMode = true // Enter stream mode
	var (
		seqNum   int64 = 0
		channels       = cl.Channels
		location       = cl.Location
		station        = text.TruncateString(cl.Station, 5)
		network        = text.TruncateString(cl.Network, 2)
	)
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
