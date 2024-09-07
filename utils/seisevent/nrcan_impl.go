package seisevent

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
	"github.com/sbabiv/xml2map"
)

const NRCAN_ID = "nrcan"

type NRCAN struct {
	cache cache.BytesCache
}

func (s *NRCAN) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      NRCAN_ID,
		Country: "CA",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Ministry of Energy and Natural Resources of Canada",
			"zh-TW": "加拿大自然資源部",
			"zh-CN": "加拿大自然资源部",
		},
	}
}

func (s *NRCAN) GetEvents(latitude, longitude float64) ([]Event, error) {
	if !s.cache.Valid() {
		res, err := request.GET(
			"https://www.earthquakescanada.nrcan.gc.ca/cache/earthquakes/canada-30.xml",
			30*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		s.cache.Set(res)
	}

	// Parse NRCAN XML response
	dataMap, err := xml2map.NewDecoder(strings.NewReader(string(s.cache.Get()))).Decode()
	if err != nil {
		return nil, err
	}
	dataMapRoot, ok := dataMap["1.2:quakeml"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("source data is not valid, missing 1.2:quakeml")
	}
	dataMapEventParameters, ok := dataMapRoot["1.2:eventParameters"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("source data is not valid, missing 1.2:eventParameters")
	}
	dataMapEvents, ok := dataMapEventParameters["1.2:event"].([]map[string]any)
	if !ok {
		return nil, errors.New("seismic event data is not available")
	}

	// Ensure the response has the expected keys and values
	expectedKeys := []string{"1.2:type", "1.2:typeCertainty", "1.2:description", "1.2:origin", "1.2:magnitude"}
	expectedKeysInDescription := []string{"1.2:text"}
	expectedKeysInOrigin := []string{"1.2:time", "1.2:latitude", "1.2:longitude", "1.2:depth"}
	expectedKeysInMagnitude := []string{"1.2:mag", "1.2:type"}

	var resultArr []Event
	for index, event := range dataMapEvents {
		if !isMapHasKeys(event, expectedKeys) {
			continue
		}

		// Check if the event is earthquake
		if event["1.2:type"].(string) != "earthquake" {
			continue
		}

		var (
			description = event["1.2:description"].(map[string]any)
			origin      = event["1.2:origin"].(map[string]any)
			magnitude   = event["1.2:magnitude"].(map[string]any)
		)
		if !isMapHasKeys(description, expectedKeysInDescription) || !isMapHasKeys(origin, expectedKeysInOrigin) || !isMapHasKeys(magnitude, expectedKeysInMagnitude) {
			continue
		}

		eventName, ok := description["1.2:text"].(string)
		if !ok {
			continue
		}
		eventTime, ok := origin["1.2:time"].(map[string]any)["1.2:value"].(string)
		if !ok {
			continue
		}
		eventLatitude, ok := origin["1.2:latitude"].(map[string]any)["1.2:value"].(string)
		if !ok {
			continue
		}
		eventLongitude, ok := origin["1.2:longitude"].(map[string]any)["1.2:value"].(string)
		if !ok {
			continue
		}
		eventDepth, ok := origin["1.2:depth"].(map[string]any)["1.2:value"].(string)
		if !ok {
			continue
		}
		eventMagnitude, ok := magnitude["1.2:mag"].(map[string]any)["1.2:value"].(string)
		if !ok {
			continue
		}
		eventMagType, ok := magnitude["1.2:type"].(string)
		if !ok {
			continue
		}

		seisEvent := Event{
			Verfied:   event["1.2:typeCertainty"].(string) == "known",
			Region:    eventName,
			Latitude:  string2Float(eventLatitude),
			Longitude: string2Float(eventLongitude),
			Depth:     string2Float(eventDepth),
			Magnitude: s.getMagnitude(eventMagType, eventMagnitude),
			Timestamp: s.getTimestamp(eventTime),
			Event:     fmt.Sprintf("NRCAN-%d", index),
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	return sortSeismicEvents(resultArr), nil
}

func (s *NRCAN) getTimestamp(data string) int64 {
	t, _ := time.Parse("2006-01-02T15:04:05.000000Z", data)
	return t.UnixMilli()
}

func (s *NRCAN) getMagnitude(magType, data string) []Magnitude {
	return []Magnitude{{
		Type:  ParseMagnitude(magType),
		Value: string2Float(data),
	}}
}
