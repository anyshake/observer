package seisevent

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const FJEA_ID = "fjea_e"

type FJEA struct {
	cache cache.BytesCache
}

func (s *FJEA) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      FJEA_ID,
		Country: "CN",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Fujian Earthquake Administration",
			"zh-TW": "福建地震局",
			"zh-CN": "福建地震局",
		},
	}
}

func (s *FJEA) GetEvents(latitude, longitude float64) ([]Event, error) {
	if !s.cache.Valid() {
		res, err := request.GET(
			"http://218.5.2.111:9088/earthquakeWarn/bulletin/list.json?pageSize=100",
			10*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		s.cache.Set(res)
	}

	// Parse FJEA JSON response
	var dataMap map[string]any
	err := json.Unmarshal(s.cache.Get(), &dataMap)
	if err != nil {
		return nil, err
	}

	dataMapEvents, ok := dataMap["list"].([]any)
	if !ok {
		return nil, errors.New("seismic event data is not available")
	}

	// Ensure the response has the expected keys and values
	expectedKeys := []string{"eventId", "shockTime", "longitude", "latitude", "placeName", "magnitude", "depth", "infoTypeName"}

	var resultArr []Event
	for _, event := range dataMapEvents {
		if !isMapHasKeys(event.(map[string]any), expectedKeys) || !isMapKeysEmpty(event.(map[string]any), expectedKeys) {
			continue
		}

		seisEvent := Event{
			Depth:     event.(map[string]any)["depth"].(float64),
			Verfied:   event.(map[string]any)["infoTypeName"].(string) == "[正式测定]",
			Event:     event.(map[string]any)["eventId"].(string),
			Region:    event.(map[string]any)["placeName"].(string),
			Latitude:  event.(map[string]any)["latitude"].(float64),
			Longitude: event.(map[string]any)["longitude"].(float64),
			Magnitude: s.getMagnitude(event.(map[string]any)["magnitude"].(float64)),
			Timestamp: s.getTimestamp(event.(map[string]any)["shockTime"].(string)),
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	return sortSeismicEvents(resultArr), nil
}

func (s *FJEA) getTimestamp(data string) int64 {
	t, _ := time.Parse("2006-01-02 15:04:05", data)
	return t.Add(-8 * time.Hour).UnixMilli()
}

func (s *FJEA) getMagnitude(data float64) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("M"), Value: data}}
}
