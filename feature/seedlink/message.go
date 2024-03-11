package seedlink

import (
	"fmt"
	"net"

	"github.com/anyshake/observer/driver/seedlink"
	"github.com/anyshake/observer/publisher"
)

func (s *SeedLink) handleMessage(conn net.Conn, client *seedlink.SeedLinkClient, gp *publisher.Geophone) error {
	if len(client.Channels) == 0 {
		return fmt.Errorf("no channels selected")
	}

	var (
		ts    = gp.TS
		ehz   = gp.EHZ
		ehe   = gp.EHE
		ehn   = gp.EHN
		chMap = map[string]publisher.Int32Array{
			"EHZ": ehz, "EHE": ehe, "EHN": ehn,
		}
	)

	for _, channel := range client.Channels {
		countData, ok := chMap[channel]
		if !ok {
			conn.Write([]byte(seedlink.RES_ERR))
			err := fmt.Errorf("channel %s not found", channel)
			return err
		}

		err := seedlink.SendSLPacket(conn, client, seedlink.SeedLinkPacket{
			Channel: channel, Timestamp: ts, Count: countData,
		})
		if err != nil {
			return err
		}
	}

	s.OnReady(nil, "SENT OK: ", conn.RemoteAddr().String())
	return nil
}
