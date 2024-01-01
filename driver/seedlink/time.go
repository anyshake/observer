package seedlink

import (
	"net"

	"github.com/anyshake/observer/feature"
)

type TIME struct{}

// Callback of "TIME <...>" command, implements SeedLinkCommandCallback interface
func (*TIME) Callback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, streamer SeedLinkStreamer, conn net.Conn, args ...string) error {
	_, err := conn.Write([]byte(RES_OK))
	return err
}

// Fallback of "TIME <...>" command, implements SeedLinkCommandCallback interface
func (*TIME) Fallback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, conn net.Conn, args ...string) {
	conn.Close()
}
