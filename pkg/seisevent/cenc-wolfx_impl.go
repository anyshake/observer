package seisevent

import (
	"encoding/json"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
)

const CENC_WOLFX_ID = "cenc_wolfx"

type CENC_WOLFX struct {
	travelTimeTable *travel.AK135
	cache           cache.AnyCache
}

func (c *CENC_WOLFX) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      CENC_WOLFX_ID,
		Country: "CN",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "China Earthquake Networks Center (Wolfx)",
			"zh-TW": "中國地震台網中心（Wolfx）",
			"zh-CN": "中国地震台网中心（Wolfx）",
		},
	}
}

func (c *CENC_WOLFX) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		return c.cache.Get().([]Event), nil
	}

	res, err := request.GET(
		"https://api.wolfx.jp/cenc_eqlist.json",
		10*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	// Parse CENC JSON response
	var dataMapEvents map[string]any
	err = json.Unmarshal(res, &dataMapEvents)
	if err != nil {
		return nil, err
	}

	// Ensure the response has the expected keys and they are not empty
	expectedKeys := []string{"EventID", "type", "time", "location", "magnitude", "depth", "latitude", "longitude"}

	var resultArr []Event
	for _, v := range dataMapEvents {
		event, ok := v.(map[string]any)
		if !ok {
			continue
		}

		if !isMapHasKeys(event, expectedKeys) || !isMapKeysEmpty(event, expectedKeys) {
			continue
		}

		timestamp, err := c.getTimestamp(event["time"].(string))
		if err != nil {
			continue
		}

		seisEvent := Event{
			Verfied:   event["type"] == "reviewed",
			Event:     event["EventID"].(string),
			Timestamp: timestamp,
			Region:    event["location"].(string),
			Depth:     string2Float(event["depth"].(string)),
			Latitude:  string2Float(event["latitude"].(string)),
			Longitude: string2Float(event["longitude"].(string)),
			Magnitude: []Magnitude{{Type: "M", Value: string2Float(event["magnitude"].(string))}},
		}

		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(c.travelTimeTable, latitude, seisEvent.Latitude, longitude, seisEvent.Longitude, seisEvent.Depth)

		resultArr = append(resultArr, seisEvent)
	}

	sortedEvents := sortSeismicEvents(resultArr)
	c.cache.Set(sortedEvents)
	return sortedEvents, nil
}

func (c *CENC_WOLFX) getTimestamp(timeStr string) (int64, error) {
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return 0, err
	}

	return t.Add(-8 * time.Hour).UnixMilli(), nil
}
