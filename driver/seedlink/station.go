package seedlink

import (
	"net"

	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/utils/text"
)

type STATION struct{}

// Callback of "STATION <...> <...>" command, implements SeedLinkCommandCallback interface
func (*STATION) Callback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, streamer SeedLinkStreamer, conn net.Conn, args ...string) error {
	cl.Station = text.TruncateString(args[0], 5)
	cl.Network = text.TruncateString(args[1], 2)
	_, err := conn.Write([]byte(RES_OK))
	return err
}

// Fallback of "STATION <...> <...>" command, implements SeedLinkCommandCallback interface
func (*STATION) Fallback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, conn net.Conn, args ...string) {
	conn.Close()
}
