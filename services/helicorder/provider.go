package helicorder

import (
	"fmt"
	"time"

	"github.com/anyshake/observer/drivers/dao/tables"
	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/utils/cache"
	"github.com/bclswl0827/heligo"
	"gorm.io/gorm"
)

type provider struct {
	database   *gorm.DB
	queryCache cache.AnyCache

	stationCode   string
	networkCode   string
	locationCode  string
	channelPrefix string

	channelCode string
}

func (d *provider) setChannelCode(channelCode string) {
	d.channelCode = channelCode
}

func (d *provider) GetPlotName() string {
	return "AnyShake Observer"
}

func (d *provider) GetStation() string {
	return d.stationCode
}

func (d *provider) GetNetwork() string {
	return d.networkCode
}

func (d *provider) GetChannel() string {
	return fmt.Sprintf("%s%s", d.channelPrefix, d.channelCode)
}

func (d *provider) GetLocation() string {
	return d.locationCode
}

func (d *provider) GetPlotData(startTime, endTime time.Time) ([]heligo.PlotData, error) {
	if !d.queryCache.Valid() {
		var (
			adcCountModel tables.AdcCount
			adcCountData  []tables.AdcCount
		)
		err := d.database.
			Table(adcCountModel.GetName()).
			Where("timestamp >= ? AND timestamp <= ?", startTime.UnixMilli(), endTime.UnixMilli()).
			Order("timestamp ASC").
			Find(&adcCountData).
			Error
		if err != nil {
			return nil, err
		}
		d.queryCache.Set(adcCountData)
	}

	var plotData []heligo.PlotData
	for _, adcCount := range d.queryCache.Get().([]tables.AdcCount) {
		data := make([]heligo.PlotData, adcCount.SampleRate)
		for i := 0; i < adcCount.SampleRate; i++ {
			timeOffset := time.Duration(i*int(time.Second.Seconds())/adcCount.SampleRate) * time.Millisecond
			data[i].Time = time.UnixMilli(adcCount.Timestamp).Add(timeOffset)
			switch d.channelCode {
			case explorer.EXPLORER_CHANNEL_CODE_Z:
				data[i].Value = float64(adcCount.Z_Axis[i])
			case explorer.EXPLORER_CHANNEL_CODE_E:
				data[i].Value = float64(adcCount.E_Axis[i])
			case explorer.EXPLORER_CHANNEL_CODE_N:
				data[i].Value = float64(adcCount.N_Axis[i])
			}
		}

		plotData = append(plotData, data...)
	}

	return plotData, nil
}
