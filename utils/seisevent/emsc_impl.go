package seisevent

import (
	"encoding/csv"
	"strconv"
	"strings"
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

type EMSC struct {
	cache cache.BytesCache
}

func (c *EMSC) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      EMSC_ID,
		Country: "EU",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "European-Mediterranean Seismological Centre",
			"zh-TW": "歐洲與地中海地震中心",
			"zh-CN": "欧洲与地中海地震中心",
		},
	}
}

func (c *EMSC) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		res, err := request.GET(
			"https://www.seismicportal.eu/fdsnws/event/1/query?minmag=-1&format=text&limit=100",
			10*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		c.cache.Set(res)
	}

	// Parse EMSC CSV response
	csvDataStr := strings.ReplaceAll(string(c.cache.Get()), ",", " - ")
	csvDataStr = strings.ReplaceAll(csvDataStr, "|", ",")
	csvRecords, err := csv.NewReader(strings.NewReader(csvDataStr)).ReadAll()
	if err != nil {
		return nil, err
	}

	var resultArr []Event
	for _, record := range csvRecords[1:] {
		var seisEvent Event
		for idx, val := range record {
			switch idx {
			case 0:
				seisEvent.Event = val
			case 1:
				seisEvent.Verfied = true
				seisEvent.Timestamp = c.getTimestamp(val)
			case 2:
				seisEvent.Latitude = c.getLatitude(val)
			case 3:
				seisEvent.Longitude = c.getLongitude(val)
			case 4:
				seisEvent.Depth = c.getDepth(val)
			case 10:
				seisEvent.Magnitude = c.getMagnitude(val)
			case 12:
				seisEvent.Region = val
			}
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	return sortSeismicEvents(resultArr), nil
}

func (c *EMSC) getTimestamp(data string) int64 {
	t, _ := time.Parse("2006-01-02T15:04:05Z", data)
	return t.UnixMilli()
}

func (c *EMSC) getMagnitude(data string) float64 {
	m, err := strconv.ParseFloat(data, 64)
	if err == nil {
		return m
	}

	return 0
}

func (c *EMSC) getDepth(data string) float64 {
	d, err := strconv.ParseFloat(data, 64)
	if err == nil {
		return d
	}

	return 0
}

func (c *EMSC) getLatitude(data string) float64 {
	lat, err := strconv.ParseFloat(data, 64)
	if err == nil {
		return lat
	}

	return 0
}

func (c *EMSC) getLongitude(data string) float64 {
	lng, err := strconv.ParseFloat(data, 64)
	if err == nil {
		return lng
	}

	return 0
}
