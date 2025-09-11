package seisevent

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
)

const GEONET_ID = "geonet"

type GEONET struct {
	travelTimeTable *travel.AK135
	cache           cache.AnyCache
}

func (u *GEONET) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      GEONET_ID,
		Country: "NZ",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "The GeoNet Project",
			"zh-TW": "GeoNet 計畫",
			"zh-CN": "GeoNet 计划",
		},
	}
}

func (u *GEONET) GetEvents(latitude, longitude float64) ([]Event, error) {
	if u.cache.Valid() {
		return u.cache.Get().([]Event), nil
	}

	res, err := request.GET(
		"https://api.geonet.org.nz/quake?MMI=1",
		10*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	// Parse GEONET JSON response
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
	expectedKeysInProperties := []string{"publicID", "time", "depth", "magnitude", "locality"}
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
		if len(coordinates.([]any)) != 2 {
			continue
		}

		seisEvent := Event{
			Verfied:   true,
			Latitude:  coordinates.([]any)[1].(float64),
			Longitude: coordinates.([]any)[0].(float64),
			Depth:     properties.(map[string]any)["depth"].(float64),
			Event:     properties.(map[string]any)["publicID"].(string),
			Region:    properties.(map[string]any)["locality"].(string),
			Magnitude: u.getMagnitude(properties.(map[string]any)["magnitude"].(float64)),
			Timestamp: u.getTimestamp(properties.(map[string]any)["time"].(string)),
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(u.travelTimeTable, latitude, seisEvent.Latitude, longitude, seisEvent.Longitude, seisEvent.Depth)

		resultArr = append(resultArr, seisEvent)
	}

	sortedEvents := sortSeismicEvents(resultArr)
	u.cache.Set(sortedEvents)
	return sortedEvents, nil
}

func (u *GEONET) getTimestamp(data string) int64 {
	t, _ := time.Parse("2006-01-02T15:04:05.000Z", data)
	return t.UnixMilli()
}

func (u *GEONET) getMagnitude(data float64) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("M"), Value: data}}
}
