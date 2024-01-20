package seedlink

import (
	"fmt"
	"net"

	"github.com/anyshake/observer/driver/seedlink"
	"github.com/anyshake/observer/publisher"
)

func (s *SeedLink) handleMessage(gp *publisher.Geophone, conn net.Conn, channels []string, network, station, location string, seqNum *int64) error {
	var (
		ehz   = gp.EHZ
		ehe   = gp.EHE
		ehn   = gp.EHN
		chMap = map[string]publisher.Int32Array{
			"EHZ": ehz, "EHE": ehe, "EHN": ehn,
		}
	)

	for _, channel := range channels {
		data, ok := chMap[channel]
		if !ok {
			conn.Write([]byte(seedlink.RES_ERR))
			err := fmt.Errorf("channel %s not found", channel)
			return err
		}

		dataBytes, err := seedlink.CreateSLPacket(data, gp.TS, *seqNum, network, station, channel, location)
		if err != nil {
			return err
		}

		_, err = conn.Write(dataBytes)
		if err != nil {
			return err
		}

		*seqNum++
	}

	s.OnReady(nil, "SENT OK: ", conn.RemoteAddr().String())
	return nil
}
