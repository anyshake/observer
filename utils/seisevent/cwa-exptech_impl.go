package seisevent

import (
	"encoding/json"
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
	"golang.org/x/exp/rand"
)

const CWA_EXPTECH_ID = "cwa_exptech"

type CWA_EXPTECH struct {
	cache cache.BytesCache
}

func (c *CWA_EXPTECH) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      CWA_EXPTECH_ID,
		Country: "TW",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Central Weather Administration (ExpTech)",
			"zh-TW": "交通部中央氣象署（ExpTech）",
			"zh-CN": "交通部中央气象署（ExpTech）",
		},
	}
}

func (c *CWA_EXPTECH) GetEvents(latitude, longitude float64) ([]Event, error) {
	// Get CWA JSON response
	if !c.cache.Valid() {
		addrs := []string{
			"https://api-1.exptech.dev/api/v2/eq/report?limit=100",
			"https://api-2.exptech.dev/api/v2/eq/report?limit=100",
		}
		res, err := request.GET(
			addrs[rand.Intn(len(addrs))],
			10*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		c.cache.Set(res)
	}

	// Parse CWA JSON response
	var dataMapEvents []map[string]any
	err := json.Unmarshal(c.cache.Get(), &dataMapEvents)
	if err != nil {
		return nil, err
	}

	// Ensure the response has the expected keys and values
	expectedKeys := []string{"id", "lat", "lon", "loc", "depth", "mag", "time"}

	var resultArr []Event
	for _, event := range dataMapEvents {
		if !isMapHasKeys(event, expectedKeys) || !isMapKeysEmpty(event, expectedKeys) {
			continue
		}

		seisEvent := Event{
			Verfied:   true,
			Depth:     event["depth"].(float64),
			Event:     event["id"].(string),
			Region:    event["loc"].(string),
			Latitude:  event["lat"].(float64),
			Longitude: event["lon"].(float64),
			Magnitude: []Magnitude{{Type: ParseMagnitude("ML"), Value: event["mag"].(float64)}},
			Timestamp: time.UnixMilli(int64(event["time"].(float64))).UnixMilli(),
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	return sortSeismicEvents(resultArr), nil
}
