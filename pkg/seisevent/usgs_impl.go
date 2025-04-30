package seisevent

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/corpix/uarand"
)

const USGS_ID = "usgs"

type USGS struct {
	cache cache.AnyCache
}

func (u *USGS) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      USGS_ID,
		Country: "US",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "United States Geological Survey",
			"zh-TW": "美國地質調查局",
			"zh-CN": "美国地质调查局",
		},
	}
}

func (u *USGS) GetEvents(latitude, longitude float64) ([]Event, error) {
	if u.cache.Valid() {
		return u.cache.Get().([]Event), nil
	}

	res, err := request.GET(
		"https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/2.5_day.geojson",
		10*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	// Parse USGS JSON response
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
	expectedKeysInDataMap := []string{"properties", "geometry", "id"}
	expectedKeysInProperties := []string{"mag", "place", "time", "type", "title", "status", "magType"}
	expectedKeysInGeometry := []string{"coordinates"}

	var resultArr []Event
	for _, event := range dataMapEvents {
		if !isMapHasKeys(event.(map[string]any), expectedKeysInDataMap) {
			continue
		}

		var (
			properties = event.(map[string]any)["properties"]
			geometry   = event.(map[string]any)["geometry"]
		)
		if !isMapHasKeys(properties.(map[string]any), expectedKeysInProperties) || !isMapHasKeys(geometry.(map[string]any), expectedKeysInGeometry) {
			continue
		}

		coordinates := geometry.(map[string]any)["coordinates"]
		if properties.(map[string]any)["type"].(string) != "earthquake" || len(coordinates.([]any)) != 3 {
			continue
		}

		seisEvent := Event{
			Depth:     coordinates.([]any)[2].(float64),
			Verfied:   properties.(map[string]any)["status"].(string) == "reviewed",
			Timestamp: int64(properties.(map[string]any)["time"].(float64)),
			Event:     event.(map[string]any)["id"].(string),
			Region:    properties.(map[string]any)["place"].(string),
			Latitude:  coordinates.([]any)[1].(float64),
			Longitude: coordinates.([]any)[0].(float64),
			Magnitude: u.getMagnitude(
				properties.(map[string]any)["magType"].(string),
				properties.(map[string]any)["mag"].(float64),
			),
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	sortedEvents := sortSeismicEvents(resultArr)
	u.cache.Set(sortedEvents)
	return sortedEvents, nil
}

func (u *USGS) getMagnitude(magType string, data float64) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude(magType), Value: data}}
}
