package seedlink

import (
	"net"

	"github.com/anyshake/observer/feature"
)

type DATA struct{}

// Callback of "DATA" command, implements SeedLinkCommandCallback interface
func (*DATA) Callback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, streamer SeedLinkStreamer, conn net.Conn, args ...string) error {
	_, err := conn.Write([]byte(RES_ERR))
	return err
}

// Fallback of "DATA" command, implements SeedLinkCommandCallback interface
func (*DATA) Fallback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, conn net.Conn, args ...string) {
	conn.Close()
}
