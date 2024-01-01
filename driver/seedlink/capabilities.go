package seedlink

import (
	"net"

	"github.com/anyshake/observer/feature"
)

type CAPABILITIES struct{}

// Callback of "CAPABILITIES" command, implements SeedLinkCommandCallback interface
func (*CAPABILITIES) Callback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, streamer SeedLinkStreamer, conn net.Conn, args ...string) error {
	_, err := conn.Write([]byte(RES_OK))
	return err
}

// Fallback of "CAPABILITIES" command, implements SeedLinkCommandCallback interface
func (*CAPABILITIES) Fallback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, conn net.Conn, args ...string) {
	conn.Close()
}
