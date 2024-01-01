package seedlink

import (
	"net"

	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/utils/text"
)

type END struct{}

// Callback of "END" command, implements SeedLinkCommandCallback interface
func (*END) Callback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, streamer SeedLinkStreamer, conn net.Conn, args ...string) error {
	cl.WorkingMode = WORKINGMODE_STREAM // Enter stream mode
	var (
		network  = text.TruncateString(options.Config.Station.Network, 2)
		station  = text.TruncateString(options.Config.Station.Station, 5)
		location = text.TruncateString(options.Config.Station.Location, 2)
	)
	go streamer(&options.Status.Geophone, conn, cl.Channel, network, station, location, &cl.WorkingMode)
	return nil
}

// Fallback of "END" command, implements SeedLinkCommandCallback interface
func (*END) Fallback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, conn net.Conn, args ...string) {
	conn.Close()
}
