package seedlink

import (
	"fmt"
	"time"

	"github.com/anyshake/observer/drivers/dao/tables"
	"github.com/anyshake/observer/utils/timesource"
	"github.com/bclswl0827/slgo/handlers"
	"gorm.io/gorm"
)

type provider struct {
	timeSource    *timesource.Source
	database      *gorm.DB
	startTime     time.Time
	stationCode   string
	networkCode   string
	locationCode  string
	channelPrefix string
}

func (p *provider) GetSoftware() string {
	return "anyshake_observer"
}

func (p *provider) GetStartTime() time.Time {
	return p.startTime
}

func (p *provider) GetCurrentTime() time.Time {
	currentTime := p.timeSource.Get()
	return currentTime
}

func (p *provider) GetOrganization() string {
	return "anyshake.org"
}

func (p *provider) GetStations() []handlers.SeedLinkStation {
	return []handlers.SeedLinkStation{
		{
			BeginSequence: "000000",
			EndSequence:   "FFFFFF",
			Station:       p.stationCode,
			Network:       p.networkCode,
			Description:   "AnyShake Observer station",
		},
	}
}

func (p *provider) GetStreams() []handlers.SeedLinkStream {
	return []handlers.SeedLinkStream{
		{
			BeginTime: p.GetStartTime().Format("2006-01-02 15:04:01"),
			EndTime:   "9999-12-31 23:59:59",
			SeedName:  fmt.Sprintf("%sZ", p.channelPrefix),
			Location:  p.locationCode,
			Type:      "D",
			Station:   p.stationCode,
		},
		{
			BeginTime: p.GetStartTime().Format("2006-01-02 15:04:01"),
			EndTime:   "9999-12-31 23:59:59",
			SeedName:  fmt.Sprintf("%sE", p.channelPrefix),
			Location:  p.locationCode,
			Type:      "D",
			Station:   p.stationCode,
		},
		{
			BeginTime: p.GetStartTime().Format("2006-01-02 15:04:01"),
			EndTime:   "9999-12-31 23:59:59",
			SeedName:  fmt.Sprintf("%sN", p.channelPrefix),
			Location:  p.locationCode,
			Type:      "D",
			Station:   p.stationCode,
		},
	}
}

func (p *provider) GetCapabilities() []handlers.SeedLinkCapability {
	return []handlers.SeedLinkCapability{
		{Name: "info:all"}, {Name: "info:gaps"}, {Name: "info:streams"},
		{Name: "dialup"}, {Name: "info:id"}, {Name: "multistation"},
		{Name: "window-extraction"}, {Name: "info:connections"},
		{Name: "info:capabilities"}, {Name: "info:stations"},
	}
}

func (p *provider) QueryHistory(startTime, endTime time.Time, channels []string) ([]handlers.SeedLinkDataPacket, error) {
	var (
		adcCountModel tables.AdcCount
		adcCountData  []tables.AdcCount
	)
	err := p.database.
		Table(adcCountModel.GetName()).
		Where("timestamp >= ? AND timestamp <= ?", startTime.UnixMilli(), endTime.UnixMilli()).
		Order("timestamp ASC").
		Find(&adcCountData).
		Error
	if err != nil {
		return nil, err
	}
	if len(adcCountData) == 0 {
		return nil, nil
	}

	// Convert ADC count data to SeedLink data packets
	var dataPacketArr []handlers.SeedLinkDataPacket
	for _, channel := range channels {
		prevSampleRate := adcCountData[0].SampleRate
		switch channel {
		case fmt.Sprintf("%sZ", p.channelPrefix):
			for _, adcCount := range adcCountData {
				if adcCount.SampleRate == prevSampleRate {
					dataPacketArr = append(dataPacketArr, handlers.SeedLinkDataPacket{
						Timestamp:  adcCount.Timestamp,
						SampleRate: adcCount.SampleRate,
						Channel:    channel,
						DataArr:    adcCount.Z_Axis,
					})
				}
				prevSampleRate = adcCount.SampleRate
			}
		case fmt.Sprintf("%sE", p.channelPrefix):
			for _, adcCount := range adcCountData {
				if adcCount.SampleRate == prevSampleRate {
					dataPacketArr = append(dataPacketArr, handlers.SeedLinkDataPacket{
						Timestamp:  adcCount.Timestamp,
						SampleRate: adcCount.SampleRate,
						Channel:    channel,
						DataArr:    adcCount.E_Axis,
					})
				}
				prevSampleRate = adcCount.SampleRate
			}
		case fmt.Sprintf("%sN", p.channelPrefix):
			for _, adcCount := range adcCountData {
				if adcCount.SampleRate == prevSampleRate {
					dataPacketArr = append(dataPacketArr, handlers.SeedLinkDataPacket{
						Timestamp:  adcCount.Timestamp,
						SampleRate: adcCount.SampleRate,
						Channel:    channel,
						DataArr:    adcCount.N_Axis,
					})
				}
			}
		}
	}

	return dataPacketArr, nil
}
