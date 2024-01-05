package seedlink

import (
	"net"

	"github.com/anyshake/observer/feature"
)

type SELECT struct{}

// Callback of "SELECT <...>" command, implements SeedLinkCommandCallback interface
func (*SELECT) Callback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, streamer SeedLinkStreamer, conn net.Conn, args ...string) error {
	if len(args) < 1 {
		_, err := conn.Write([]byte(RES_ERR))
		return err
	} else {
		if len(args[0]) < 5 {
			_, err := conn.Write([]byte(RES_ERR))
			return err
		} else {
			cl.Location = args[0][:2]
			cl.Channels = append(cl.Channels, args[0][2:5])
		}
	}
	_, err := conn.Write([]byte(RES_OK))
	return err
}

// Fallback of "SELECT <...>" command, implements SeedLinkCommandCallback interface
func (*SELECT) Fallback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, conn net.Conn, args ...string) {
	conn.Close()
}
