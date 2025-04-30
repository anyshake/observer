package seisevent

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/corpix/uarand"
)

const SCEA_E_ID = "scea_e"

type SCEA_E struct {
	cache cache.AnyCache
}

func (s *SCEA_E) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      SCEA_E_ID,
		Country: "CN",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Sichuan Earthquake Administration (Early Warning)",
			"zh-TW": "四川地震局（預警）",
			"zh-CN": "四川地震局（预警）",
		},
	}
}

func (s *SCEA_E) GetEvents(latitude, longitude float64) ([]Event, error) {
	if s.cache.Valid() {
		return s.cache.Get().([]Event), nil
	}

	res, err := request.GET(
		"http://118.113.105.29:8002/api/earlywarning/jsonPageList?pageSize=100",
		10*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	// Parse SCEA_B JSON response
	var dataMap map[string]any
	err = json.Unmarshal(res, &dataMap)
	if err != nil {
		return nil, err
	}

	// Check server response
	if dataMap["code"].(float64) != 0 {
		return nil, fmt.Errorf("server error: %s", dataMap["msg"])
	}

	dataMapEvents, ok := dataMap["data"].([]any)
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
			Depth:     -1,
			Verfied:   event.(map[string]any)["infoTypeName"].(string) == "[正式]",
			Event:     event.(map[string]any)["eventId"].(string),
			Region:    event.(map[string]any)["placeName"].(string),
			Latitude:  event.(map[string]any)["latitude"].(float64),
			Longitude: event.(map[string]any)["longitude"].(float64),
			Magnitude: s.getMagnitude(event.(map[string]any)["magnitude"].(float64)),
			Timestamp: time.UnixMilli(int64(event.(map[string]any)["shockTime"].(float64))).UnixMilli(),
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	sortedEvents := sortSeismicEvents(resultArr)
	s.cache.Set(sortedEvents)
	return sortedEvents, nil
}

func (s *SCEA_E) getMagnitude(data float64) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("M"), Value: data}}
}
