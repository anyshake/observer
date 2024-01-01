package seedlink

import (
	"net"

	"github.com/anyshake/observer/feature"
)

type STATION struct{}

// Callback of "STATION <...> <...>" command, implements SeedLinkCommandCallback interface
func (*STATION) Callback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, streamer SeedLinkStreamer, conn net.Conn, args ...string) error {
	if len(args) < 3 {
		_, err := conn.Write([]byte(RES_ERR))
		return err
	} else {
		cl.Station = args[1]
		cl.Network = args[2]
	}

	_, err := conn.Write([]byte(RES_OK))
	return err
}

// Fallback of "STATION <...> <...>" command, implements SeedLinkCommandCallback interface
func (*STATION) Fallback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, conn net.Conn, args ...string) {
	conn.Close()
}
