package seisevent

import (
	"encoding/json"
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

type CEIC struct {
	cache cache.BytesCache
}

func (c *CEIC) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      CEIC_ID,
		Country: "CN",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "China Earthquake Networks Center",
			"zh-TW": "中國地震台網中心",
			"zh-CN": "中国地震台网中心",
		},
	}
}

func (c *CEIC) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		res, err := request.GET(
			"https://news.ceic.ac.cn/ajax/google",
			10*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		c.cache.Set(res)
	}

	// Parse CEIC JSON response
	responseMap := make([]map[string]any, 0)
	err := json.Unmarshal(c.cache.Get(), &responseMap)
	if err != nil {
		return nil, err
	}

	// Ensure the response has the expected keys and they are not empty
	expectedKeys := []string{"CATA_ID", "O_TIME", "EPI_LAT", "EPI_LON", "EPI_DEPTH", "M", "LOCATION_C"}

	var resultArr []Event
	for _, v := range responseMap {
		if !isMapHasKeys(v, expectedKeys) || !isMapKeysEmpty(v, expectedKeys) {
			continue
		}

		timestamp, err := c.getTimestamp(v["O_TIME"].(string))
		if err != nil {
			continue
		}

		seisEvent := Event{
			Verfied:   true,
			Timestamp: timestamp,
			Depth:     c.getDepth(v["EPI_DEPTH"]),
			Event:     v["CATA_ID"].(string),
			Region:    v["LOCATION_C"].(string),
			Latitude:  string2Float(v["EPI_LAT"].(string)),
			Longitude: string2Float(v["EPI_LON"].(string)),
			Magnitude: string2Float(v["M"].(string)),
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	return sortSeismicEvents(resultArr), nil
}

func (c *CEIC) getTimestamp(timeStr string) (int64, error) {
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return 0, err
	}

	return t.Add(-8 * time.Hour).UnixMilli(), nil
}

func (c *CEIC) getDepth(depth any) float64 {
	switch d := depth.(type) {
	case string:
		return string2Float(d)
	case float64:
		return d
	}

	return -1
}
