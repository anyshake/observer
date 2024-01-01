package seedlink

import (
	"fmt"
	"net"

	"github.com/anyshake/observer/feature"
)

type HELLO struct{}

// Callback of "HELLO" command, implements SeedLinkCommandCallback interface
func (*HELLO) Callback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, streamer SeedLinkStreamer, conn net.Conn, args ...string) error {
	station := sl.Organization
	_, err := conn.Write([]byte(fmt.Sprintf("%s\r\n%s\r\n", RELEASE, station)))
	return err
}

// Fallback of "HELLO" command, implements SeedLinkCommandCallback interface
func (*HELLO) Fallback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, conn net.Conn, args ...string) {
	conn.Close()
}
