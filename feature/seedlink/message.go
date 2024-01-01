package seedlink

import (
	"fmt"
	"net"
	"time"

	"github.com/anyshake/observer/driver/seedlink"
	"github.com/anyshake/observer/publisher"
	"github.com/bclswl0827/mseedio"
)

func (s *SeedLink) handleMessage(gp *publisher.Geophone, conn net.Conn, channel, network, station, location string, seqNum *int64) error {
	var (
		ehz   = gp.EHZ
		ehe   = gp.EHE
		ehn   = gp.EHN
		chMap = map[string]publisher.Int32Array{
			"EHZ": ehz, "EHE": ehe, "EHN": ehn,
		}
	)

	// Check if channel exists
	data, ok := chMap[channel]
	if !ok {
		conn.Write([]byte(seedlink.RES_ERR))
		err := fmt.Errorf("channel %s not found", channel)
		return err
	}

	// Generate MiniSEED, send to client
	var miniseed mseedio.MiniSeedData
	// Init header fields
	miniseed.Init(mseedio.STEIM2, mseedio.MSBFIRST)

	// Append MiniSEED data
	err := miniseed.Append(data, &mseedio.AppendOptions{
		StationCode:    station,
		LocationCode:   location,
		ChannelCode:    channel,
		NetworkCode:    network,
		SampleRate:     float64(len(data) - 1),
		StartTime:      time.UnixMilli(gp.TS).UTC(),
		SequenceNumber: fmt.Sprintf("%06d", *seqNum),
	})
	if err != nil {
		return err
	}

	// Get MiniSEED bytes in fixed 512 bytes
	miniseed.Series[0].BlocketteSection.RecordLength = 9
	dataBytes, err := miniseed.Encode(mseedio.OVERWRITE, mseedio.MSBFIRST)
	if err != nil {
		return err
	}

	// Send data with SL sequence number
	dataBytes = append([]byte(fmt.Sprintf("SL%06X", *seqNum)), dataBytes...)
	_, err = conn.Write(dataBytes)
	if err != nil {
		return err
	} else {
		*seqNum++
	}

	s.OnReady(nil, "SENT OK: ", conn.RemoteAddr().String())
	return nil
}
