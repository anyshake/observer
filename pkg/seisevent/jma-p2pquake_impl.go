package seisevent

import (
	"context"
	"encoding/json"
	"slices"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
	"golang.org/x/sync/singleflight"
)

const JMA_P2PQUAKE_ID = "jma_p2pquake"

type JMA_P2PQUAKE struct {
	travelTimeTable *travel.AK135
	cache           cache.GenericCache[[]Event]
	sf              singleflight.Group
}

func (j *JMA_P2PQUAKE) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      JMA_P2PQUAKE_ID,
		Country: "JP",
		Default: "en-US",
		Locales: map[string]string{
			"en-US": "Japan Meteorological Agency (P2PQuake)",
			"zh-TW": "氣象廳（P2P 地震情報）",
		},
	}
}

func (j *JMA_P2PQUAKE) GetEvents(ctx context.Context, latitude, longitude float64) ([]Event, error) {
	var baseEvents []Event

	if j.cache.Valid() {
		baseEvents = j.cache.Get()
	} else {
		v, err := fetchEvents(ctx, &j.sf, JMA_P2PQUAKE_ID, func(ctx context.Context) ([]Event, error) {
			if j.cache.Valid() {
				return j.cache.Get(), nil
			}

			res, err := request.GET(
				ctx, "https://api.p2pquake.net/v2/jma/quake",
				10*time.Second, time.Second, 3, false, nil,
				map[string]string{"User-Agent": uarand.GetRandom()},
			)
			if err != nil {
				return nil, err
			}

			// Parse JMA_P2PQUAKE JSON response
			var dataMapEvents []map[string]any
			err = json.Unmarshal(res, &dataMapEvents)
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

				resultArr = append(resultArr, Event{
					Verfied:   true,
					Timestamp: timestamp,
					Depth:     hypocenter["depth"].(float64),
					Event:     event["id"].(string),
					Region:    hypocenter["name"].(string),
					Latitude:  hypocenter["latitude"].(float64),
					Longitude: hypocenter["longitude"].(float64),
					Magnitude: []Magnitude{{Type: ParseMagnitude("M"), Value: hypocenter["magnitude"].(float64)}},
				})
			}

			sortedEvents := sortSeismicEvents(resultArr)
			j.cache.Set(sortedEvents)
			return sortedEvents, nil
		})
		if err != nil {
			return nil, err
		}

		baseEvents = v
	}

	baseEvents = slices.Clone(baseEvents)
	for i := range baseEvents {
		baseEvents[i].Distance = getDistance(latitude, baseEvents[i].Latitude, longitude, baseEvents[i].Longitude)
		baseEvents[i].Estimation = getSeismicEstimation(
			j.travelTimeTable,
			latitude,
			baseEvents[i].Latitude,
			longitude,
			baseEvents[i].Longitude,
			baseEvents[i].Depth,
		)
	}

	return baseEvents, nil
}

func (j *JMA_P2PQUAKE) getTimestamp(timeStr string) (int64, error) {
	t, err := time.Parse("2006/01/02 15:04:05", timeStr)
	if err != nil {
		return 0, err
	}

	return t.Add(-9 * time.Hour).UnixMilli(), nil
}
