package seedlink

import (
	"net"
	"time"

	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/publisher"
	"github.com/ostafen/clover/v2/query"
)

type END struct{}

// Callback of "END" command, implements SeedLinkCommandCallback interface
func (*END) Callback(sl *SeedLinkGlobal, client *SeedLinkClient, options *feature.FeatureOptions, streamer SeedLinkStreamer, conn net.Conn, args ...string) error {
	if client.StartTime.IsZero() {
		_, err := conn.Write([]byte(RES_ERR))
		return err
	}

	// Query from buffer database
	records, err := sl.SeedLinkBuffer.Database.FindAll(query.NewQuery(sl.SeedLinkBuffer.Collection).
		Where(query.Field("ts").Gt(client.StartTime.UTC().UnixMilli()).
			And(query.Field("ts").Lt(client.EndTime.UTC().UnixMilli())),
		))
	if err != nil {
		conn.Write([]byte(RES_ERR))
		return err
	}

	// Enter stream mode
	client.Streaming = true

	for _, record := range records {
		var recordMap map[string]any
		record.Unmarshal(&recordMap)
		channelMap := map[string]string{
			"EHZ": recordMap["ehz"].(string),
			"EHE": recordMap["ehe"].(string),
			"EHN": recordMap["ehn"].(string),
		}
		for _, channel := range client.Channels {
			data, ok := channelMap[channel]
			if !ok {
				continue
			}
			var (
				timestamp = int64(recordMap["ts"].(float64))
				bufTime   = time.UnixMilli(timestamp).UTC()
			)
			if bufTime.After(client.StartTime.UTC()) && bufTime.Before(client.EndTime.UTC()) {
				countDataArr, err := publisher.DecodeInt32Array(data)
				if err != nil {
					return err
				}
				err = SendSLPacket(conn, client, SeedLinkPacket{
					Channel: channel, Timestamp: timestamp, Count: countDataArr,
				})
				if err != nil {
					return err
				}
			}
		}
	}

	// Subscribe to the publisher
	go publisher.Subscribe(
		&options.Status.Geophone, &client.Streaming,
		func(gp *publisher.Geophone) error {
			return streamer(conn, client, gp)
		},
	)

	return nil
}

// Fallback of "END" command, implements SeedLinkCommandCallback interface
func (*END) Fallback(sl *SeedLinkGlobal, client *SeedLinkClient, options *feature.FeatureOptions, conn net.Conn, args ...string) {
	conn.Close()
}
