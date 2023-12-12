package seedlink

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/anyshake/observer/feature"
)

func (s *SeedLink) handleConnection(options *feature.FeatureOptions, conn net.Conn) {
	station := options.Config.SeedLink.Station
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := strings.ToUpper(scanner.Text())
		switch msg {
		case "HELLO":
			conn.Write([]byte(fmt.Sprintf("%s\r\n%s\r\n", RELEASE, station)))
		default:
			fmt.Println("Unknown message:", msg)
			conn.Write([]byte("ERROR\n"))
		}
	}

	defer conn.Close()
	s.OnStop(nil, "1 client has disconnected")
}
