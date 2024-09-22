package seisevent

import (
	"encoding/json"
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const JMA_P2PQUAKE_ID = "jma_p2pquake"

type JMA_P2PQUAKE struct {
	cache cache.BytesCache
}

func (j *JMA_P2PQUAKE) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      JMA_P2PQUAKE_ID,
		Country: "JP",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Japan Meteorological Agency (P2PQuake)",
			"zh-TW": "氣象廳（P2P 地震情報）",
			"zh-CN": "气象厅（P2P 地震情报）",
		},
	}
}

func (j *JMA_P2PQUAKE) GetEvents(latitude, longitude float64) ([]Event, error) {
	if !j.cache.Valid() {
		res, err := request.GET(
			"https://api.p2pquake.net/v2/jma/quake",
			10*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		j.cache.Set(res)
	}

	// Parse JMA_P2PQUAKE JSON response
	var dataMapEvents []map[string]any
	err := json.Unmarshal(j.cache.Get(), &dataMapEvents)
	if err != nil {
		return nil, err
	}

	// Ensure the response has the expected keys and values
	expectedKeysInDataMap := []string{"id", "earthquake"}
	expectedKeysInEarthquake := []string{"time", "hypocenter"}
	expectedKeysInHypocenter := []string{"name", "latitude", "longitude", "depth", "magnitude"}

	var resultArr []Event
	for _, event := range dataMapEvents {
		if !isMapHasKeys(event, expectedKeysInDataMap) || !isMapKeysEmpty(event, expectedKeysInDataMap) {
			continue
		}

		var (
			earthquake = event["earthquake"].(map[string]any)
			hypocenter = earthquake["hypocenter"].(map[string]any)
		)
		if !isMapHasKeys(earthquake, expectedKeysInEarthquake) || !isMapKeysEmpty(earthquake, expectedKeysInEarthquake) {
			continue
		}
		if !isMapHasKeys(hypocenter, expectedKeysInHypocenter) || !isMapKeysEmpty(hypocenter, expectedKeysInHypocenter) {
			continue
		}

		timestamp, err := j.getTimestamp(earthquake["time"].(string))
		if err != nil {
			continue
		}

		seisEvent := Event{
			Verfied:   true,
			Timestamp: timestamp,
			Depth:     hypocenter["depth"].(float64),
			Event:     event["id"].(string),
			Region:    hypocenter["name"].(string),
			Latitude:  hypocenter["latitude"].(float64),
			Longitude: hypocenter["longitude"].(float64),
			Magnitude: []Magnitude{{Type: ParseMagnitude("M"), Value: hypocenter["magnitude"].(float64)}},
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	return sortSeismicEvents(resultArr), nil
}

func (j *JMA_P2PQUAKE) getTimestamp(timeStr string) (int64, error) {
	t, err := time.Parse("2006/01/02 15:04:05", timeStr)
	if err != nil {
		return 0, err
	}

	return t.Add(-9 * time.Hour).UnixMilli(), nil
}
