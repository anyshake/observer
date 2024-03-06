package seedlink

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/publisher"
	"github.com/anyshake/observer/utils/text"
	"github.com/ostafen/clover/v2/query"
)

type END struct{}

// Callback of "END" command, implements SeedLinkCommandCallback interface
func (*END) Callback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, streamer SeedLinkStreamer, conn net.Conn, args ...string) error {
	cl.StreamMode = true // Enter stream mode
	var (
		seqNum     int64 = 0
		channels         = cl.Channels
		location         = cl.Location
		endTime          = cl.EndTime
		startTime        = cl.StartTime
		database         = sl.SeedLinkBuffer.Database
		collection       = sl.SeedLinkBuffer.Collection
		station          = text.TruncateString(cl.Station, 5)
		network          = text.TruncateString(cl.Network, 2)
	)

	if startTime.IsZero() {
		_, err := conn.Write([]byte(RES_ERR))
		return err
	}

	records, err := database.FindAll(
		query.NewQuery(collection).
			Where(query.Field("ts").
				Gt(startTime.UnixMilli()).
				And(query.Field("ts").
					Lt(endTime.UnixMilli()),
				),
			),
	)
	if err != nil {
		conn.Write([]byte(RES_ERR))
		return err
	}

	for _, record := range records {
		var recordMap map[string]any
		record.Unmarshal(&recordMap)
		channelMap := map[string]string{
			"EHZ": recordMap["ehz"].(string),
			"EHE": recordMap["ehe"].(string),
			"EHN": recordMap["ehn"].(string),
		}
		for _, channel := range channels {
			data, ok := channelMap[channel]
			if !ok {
				continue
			}
			var (
				timestamp = int64(recordMap["ts"].(float64))
				bufTime   = time.UnixMilli(timestamp)
			)
			if bufTime.After(startTime) && bufTime.Before(endTime) {
				var countData []int32
				for _, v := range strings.Split(data, "|") {
					intData, err := strconv.Atoi(v)
					if err != nil {
						return err
					}
					countData = append(countData, int32(intData))
				}
				err := SendSLPacket(conn, countData, timestamp, &seqNum, network, station, channel, location)
				if err != nil {
					return err
				}
			}
		}
	}

	// Subscribe to the publisher
	go publisher.Subscribe(
		&options.Status.Geophone, &cl.StreamMode,
		func(gp *publisher.Geophone) error {
			return streamer(gp, conn, channels, network, station, location, &seqNum)
		},
	)

	return nil
}

// Fallback of "END" command, implements SeedLinkCommandCallback interface
func (*END) Fallback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, conn net.Conn, args ...string) {
	conn.Close()
}
