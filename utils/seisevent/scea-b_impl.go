package seisevent

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const SCEA_B_ID = "scea_b"

type SCEA_B struct {
	cache cache.AnyCache
}

func (s *SCEA_B) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      SCEA_B_ID,
		Country: "CN",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Sichuan Earthquake Administration (Bulletin)",
			"zh-TW": "四川地震局（速報）",
			"zh-CN": "四川地震局（速报）",
		},
	}
}

func (s *SCEA_B) GetEvents(latitude, longitude float64) ([]Event, error) {
	if !s.cache.Valid() {
		res, err := request.GET(
			"http://118.113.105.29:8002/api/bulletin/jsonPageList?pageSize=100",
			10*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		s.cache.Set(res)
	}

	// Parse SCEA_B JSON response
	var dataMap map[string]any
	err := json.Unmarshal(s.cache.Get().([]byte), &dataMap)
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
	expectedKeys := []string{"autoFlag", "eventId", "shockTime", "longitude", "latitude", "placeName", "magnitude", "depth"}

	var resultArr []Event
	for _, event := range dataMapEvents {
		if !isMapHasKeys(event.(map[string]any), expectedKeys) || !isMapKeysEmpty(event.(map[string]any), expectedKeys) {
			continue
		}

		region := event.(map[string]any)["placeName"].(string)
		if strings.HasPrefix(region, "中国") {
			region = strings.ReplaceAll(region, "中国", "")
		}
		if strings.HasPrefix(region, "台湾省") {
			region = strings.ReplaceAll(region, "台湾", "")
		}

		seisEvent := Event{
			Region:    region,
			Verfied:   event.(map[string]any)["autoFlag"].(string) != "AU",
			Depth:     event.(map[string]any)["depth"].(float64),
			Event:     event.(map[string]any)["eventId"].(string),
			Latitude:  event.(map[string]any)["latitude"].(float64),
			Longitude: event.(map[string]any)["longitude"].(float64),
			Magnitude: s.getMagnitude(event.(map[string]any)["magnitude"].(float64)),
			Timestamp: time.UnixMilli(int64(event.(map[string]any)["shockTime"].(float64))).UnixMilli(),
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	return sortSeismicEvents(resultArr), nil
}

func (s *SCEA_B) getMagnitude(data float64) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("M"), Value: data}}
}
