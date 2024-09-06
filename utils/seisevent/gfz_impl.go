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

type GFZ struct {
	cache cache.BytesCache
}

func (c *GFZ) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      GFZ_ID,
		Country: "DE",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "GFZ German Research Centre",
			"zh-TW": "亥姆霍茲德國地理研究中心",
			"zh-CN": "德国亥姆霍兹地球科学研究中心",
		},
	}
}

func (c *GFZ) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		res, err := request.GET(
			"https://geofon.gfz-potsdam.de/fdsnws/event/1/query?minmag=-1&format=text&limit=100",
			10*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		c.cache.Set(res)
	}

	// Parse GFZ CSV response
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

func (c *GFZ) getTimestamp(data string) int64 {
	t, _ := time.Parse("2006-01-02T15:04:05.99", data)
	return t.UnixMilli()
}

func (c *GFZ) getMagnitude(data string) float64 {
	m, err := strconv.ParseFloat(data, 64)
	if err == nil {
		return m
	}

	return 0
}

func (c *GFZ) getDepth(data string) float64 {
	d, err := strconv.ParseFloat(data, 64)
	if err == nil {
		return d
	}

	return 0
}

func (c *GFZ) getLatitude(data string) float64 {
	lat, err := strconv.ParseFloat(data, 64)
	if err == nil {
		return lat
	}

	return 0
}

func (c *GFZ) getLongitude(data string) float64 {
	lng, err := strconv.ParseFloat(data, 64)
	if err == nil {
		return lng
	}

	return 0
}
