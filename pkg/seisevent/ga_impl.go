package seisevent

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/corpix/uarand"
)

const GA_ID = "ga"

type GA struct {
	cache cache.AnyCache
}

func (c *GA) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      GA_ID,
		Country: "AU",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Geoscience Australia",
			"zh-TW": "澳洲地球科學局",
			"zh-CN": "澳洲地球科学局",
		},
	}
}

func (c *GA) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		return c.cache.Get().([]Event), nil
	}

	res, err := request.GET(
		"https://earthquakes.ga.gov.au/cache/recent-earthquakes.json",
		30*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	// Parse JMA JSON response
	var dataMap map[string]any
	err = json.Unmarshal(res, &dataMap)
	if err != nil {
		return nil, err
	}

	dataMapEvents, ok := dataMap["features"].([]any)
	if !ok {
		return nil, errors.New("seismic event data is not available")
	}

	// Ensure the response has the expected keys and values
	expectedKeysInDataMap := []string{"properties", "geometry"}
	expectedKeysInProperties := []string{"depth", "preferred_magnitude_type", "preferred_magnitude", "event_id", "latitude", "longitude", "description", "evaluation_status", "origin_time"}

	var resultArr []Event
	for _, event := range dataMapEvents {
		if !isMapHasKeys(event.(map[string]any), expectedKeysInDataMap) {
			continue
		}

		properties := event.(map[string]any)["properties"]
		if !isMapHasKeys(properties.(map[string]any), expectedKeysInProperties) {
			continue
		}

		seisEvent := Event{
			Depth:     properties.(map[string]any)["depth"].(float64),
			Verfied:   properties.(map[string]any)["evaluation_status"].(string) == "confirmed",
			Timestamp: c.getTimestamp(properties.(map[string]any)["origin_time"].(string)),
			Event:     properties.(map[string]any)["event_id"].(string),
			Region:    properties.(map[string]any)["description"].(string),
			Latitude:  properties.(map[string]any)["latitude"].(float64),
			Longitude: properties.(map[string]any)["longitude"].(float64),
			Magnitude: c.getMagnitude(
				properties.(map[string]any)["preferred_magnitude_type"].(string),
				properties.(map[string]any)["preferred_magnitude"].(float64),
			),
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	sortedEvents := sortSeismicEvents(resultArr)
	c.cache.Set(sortedEvents)
	return sortedEvents, nil
}

func (c *GA) getTimestamp(data string) int64 {
	t, _ := time.Parse("2006-01-02T15:04:05.000Z", data)
	return t.UnixMilli()
}

func (c *GA) getMagnitude(magType string, data float64) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude(magType), Value: data}}
}
