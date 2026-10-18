package fdsnws

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/anyshake/observer/internal/dao/model"
	fdsnwsgo "github.com/bclswl0827/fdsnws-go"
)

type dataSelectProvider struct {
	*provider
	ctx context.Context
}

func (p *dataSelectProvider) DataSelect(query fdsnwsgo.DataSelectQuery) (fdsnwsgo.DataSelectResult, fdsnwsgo.ProviderErrorKind, error) {
	result := fdsnwsgo.DataSelectResult{}
	if len(query.Selections) > 100 {
		return result, fdsnwsgo.ProviderErrorClient, errors.New("at most 100 selections are allowed")
	}
	window := p.handler.SeisRecordsGetQueryWindow()
	var total time.Duration
	for _, selection := range query.Selections {
		if selection.Start == nil || selection.End == nil {
			return result, fdsnwsgo.ProviderErrorClient, errors.New("starttime and endtime are required")
		}
		duration := selection.End.Sub(*selection.Start)
		if duration < 0 || duration > window-total {
			return result, fdsnwsgo.ProviderErrorClient, fmt.Errorf("total query duration exceeds %s limit", window)
		}
		total += duration
	}
	network, station, location, err := p.stationCodes()
	if err != nil {
		return result, fdsnwsgo.ProviderErrorServer, err
	}
	// All stored observations are raw data (quality D).
	if query.Quality != "B" && query.Quality != "D" {
		return result, fdsnwsgo.ProviderErrorNone, fdsnwsgo.ErrNoData
	}
	deviceConfig := p.device.GetConfig()
	for _, selection := range query.Selections {
		if !matches(selection.Network, network) || !matches(selection.Station, station) || !matches(selection.Location, location) {
			continue
		}
		// Include the packet preceding starttime, then trim to exact sample boundaries.
		padding := max(deviceConfig.GetPacketInterval(), time.Second)
		start, end := selection.Start.Add(-padding), *selection.End
		for !start.After(end) {
			chunkEnd := start.Add(window)
			if chunkEnd.After(end) {
				chunkEnd = end
			}
			err := p.handler.SeisRecordsQueryEachContext(p.ctx, start, chunkEnd, func(record model.SeisRecord) error {
				timestamp, rate, channels, err := record.Decode()
				if err != nil {
					return err
				}
				if rate <= 0 || timestamp.After(end) {
					return nil
				}
				for _, channel := range channels {
					if !matches(selection.Channel, channel.ChannelCode) {
						continue
					}
					first := 0
					if timestamp.Before(*selection.Start) {
						first = int((selection.Start.Sub(timestamp)*time.Duration(rate) + time.Second - 1) / time.Second)
					}
					last := min(len(channel.Data), int(end.Sub(timestamp)*time.Duration(rate)/time.Second)+1)
					if first >= last {
						continue
					}
					result.Records = append(result.Records, fdsnwsgo.DataRecord{
						Network: network, Station: station, Location: location, Channel: channel.ChannelCode,
						Quality: "D", StartTime: timestamp.Add(time.Duration(first) * time.Second / time.Duration(rate)).UTC(),
						SampleRate: float64(rate), Samples: channel.Data[first:last],
					})
				}
				return nil
			})
			if err != nil {
				return fdsnwsgo.DataSelectResult{}, fdsnwsgo.ProviderErrorServer, err
			}
			start = chunkEnd.Add(time.Millisecond)
		}
	}
	result.Records = mergeRecords(result.Records)
	if len(result.Records) == 0 {
		return result, fdsnwsgo.ProviderErrorNone, fdsnwsgo.ErrNoData
	}
	return result, fdsnwsgo.ProviderErrorNone, nil
}

// Merge contiguous packets and overlapping POST selections before applying
// fdsnws-go's minimumlength/longestonly filters.
func mergeRecords(records []fdsnwsgo.DataRecord) []fdsnwsgo.DataRecord {
	sort.SliceStable(records, func(i, j int) bool {
		if records[i].Channel != records[j].Channel {
			return records[i].Channel < records[j].Channel
		}
		return records[i].StartTime.Before(records[j].StartTime)
	})
	merged := make([]fdsnwsgo.DataRecord, 0, len(records))
	for _, record := range records {
		if len(merged) > 0 {
			last := &merged[len(merged)-1]
			if last.Channel == record.Channel && last.SampleRate == record.SampleRate {
				rate := int64(last.SampleRate)
				next := last.StartTime.Add(time.Duration(len(last.Samples)) * time.Second / time.Duration(rate))
				if !record.StartTime.After(next) {
					skip := int((next.Sub(record.StartTime)*time.Duration(rate) + time.Second - 1) / time.Second)
					if skip < len(record.Samples) {
						last.Samples = append(last.Samples, record.Samples[skip:]...)
					}
					continue
				}
			}
		}
		record.Samples = append([]int32(nil), record.Samples...)
		merged = append(merged, record)
	}
	return merged
}
