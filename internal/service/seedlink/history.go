package seedlink

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/bclswl0827/slgo/handlers"
)

func (p *provider) QueryHistory(startTime, endTime time.Time, channels []handlers.SeedLinkChannel) ([]handlers.SeedLinkDataPacket, error) {
	if endTime.IsZero() {
		endTime = p.timeSource.Now()
	}

	channelSet := make(map[string]struct{}, len(channels))
	for _, ch := range channels {
		channelSet[ch.ChannelName] = struct{}{}
	}
	if len(channelSet) == 0 {
		return nil, nil
	}
	if startTime.After(endTime) {
		return nil, errors.New("start time is after end time")
	}
	queryWindow := p.actionHandler.SeisRecordsGetQueryWindow()
	if endTime.Sub(startTime) > queryWindow {
		return nil, fmt.Errorf("duration between start time and end time exceeds %.0f minutes limit", queryWindow.Minutes())
	}

	type packetKey struct {
		timestamp int64
		channel   string
	}
	packets := make(map[packetKey]handlers.SeedLinkDataPacket)
	addPacket := func(packet handlers.SeedLinkDataPacket) {
		if _, exists := channelSet[packet.Channel]; !exists || packet.SampleRate <= 0 || len(packet.DataArr) == 0 {
			return
		}
		recordStart := time.UnixMilli(packet.Timestamp)
		recordDuration := time.Duration(len(packet.DataArr)/packet.SampleRate)*time.Second +
			time.Duration(len(packet.DataArr)%packet.SampleRate)*time.Second/time.Duration(packet.SampleRate)
		if recordStart.Before(endTime) && recordStart.Add(recordDuration).After(startTime) {
			packets[packetKey{timestamp: packet.Timestamp, channel: packet.Channel}] = packet
		}
	}

	// Retain this snapshot across the database query: a slow query must not
	// lose uncommitted records when the shared ring advances in the meantime.
	var buffered []explorer.Event
	if p.recentBuffer != nil {
		buffered = p.recentBuffer.Values()
	}

	collect := func(record model.SeisRecord) error {
		_, sampleRate, channelData, err := record.Decode()
		if err != nil {
			return err
		}
		for _, data := range channelData {
			addPacket(handlers.SeedLinkDataPacket{
				Timestamp: record.RecordTime, SampleRate: sampleRate,
				Channel: data.ChannelCode, DataArr: data.Data,
			})
		}
		return nil
	}

	// Hardware archive events contain one second of samples. Include the event
	// crossing the requested start; slgo trims it on the sample grid. Split at
	// the DAO limit so an otherwise valid one-hour request still succeeds.
	for queryStart := startTime.Add(-time.Second); queryStart.Before(endTime); {
		queryEnd := endTime
		if queryEnd.Sub(queryStart) > queryWindow {
			queryEnd = queryStart.Add(queryWindow)
		}
		if err := p.actionHandler.SeisRecordsQueryEach(p.ctx, queryStart, queryEnd, collect); err != nil {
			return nil, err
		}
		queryStart = queryEnd
	}

	if p.recentBuffer != nil {
		buffered = append(buffered, p.recentBuffer.Values()...)
	}
	// Buffered data wins over database duplicates and covers records that the
	// archiver had not committed when the query ran.
	for _, event := range buffered {
		for _, data := range event.ChannelData {
			addPacket(handlers.SeedLinkDataPacket{
				Timestamp: event.Timestamp.UnixMilli(), SampleRate: event.SampleRate,
				Channel: data.ChannelCode, DataArr: data.Data,
			})
		}
	}

	dataPackets := make([]handlers.SeedLinkDataPacket, 0, len(packets))
	for _, packet := range packets {
		dataPackets = append(dataPackets, packet)
	}
	sort.Slice(dataPackets, func(i, j int) bool {
		if dataPackets[i].Timestamp == dataPackets[j].Timestamp {
			return dataPackets[i].Channel < dataPackets[j].Channel
		}
		return dataPackets[i].Timestamp < dataPackets[j].Timestamp
	})
	return dataPackets, nil
}
