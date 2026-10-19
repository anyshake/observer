package seisevent

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
	"golang.org/x/sync/singleflight"
)

const USGS_ID = "usgs"

type USGS struct {
	travelTimeTable *travel.AK135
	cache           cache.GenericCache[[]Event]
	sf              singleflight.Group
}

func (u *USGS) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      USGS_ID,
		Country: "US",
		Default: "en-US",
		Locales: map[string]string{
			"en-US": "United States Geological Survey",
			"zh-TW": "美國地質調查局",
		},
	}
}

func (u *USGS) GetEvents(ctx context.Context, latitude, longitude float64) ([]Event, error) {
	var baseEvents []Event

	if u.cache.Valid() {
		baseEvents = u.cache.Get()
	} else {
		v, err := fetchEvents(ctx, &u.sf, USGS_ID, func(ctx context.Context) ([]Event, error) {
			if u.cache.Valid() {
				return u.cache.Get(), nil
			}

			res, err := request.GET(
				ctx, "https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/2.5_day.geojson",
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

				resultArr = append(resultArr, Event{
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
				})
			}

			sortedEvents := sortSeismicEvents(resultArr)
			u.cache.Set(sortedEvents)
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
			u.travelTimeTable,
			latitude,
			baseEvents[i].Latitude,
			longitude,
			baseEvents[i].Longitude,
			baseEvents[i].Depth,
		)
	}

	return baseEvents, nil
}

func (u *USGS) getMagnitude(magType string, data float64) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude(magType), Value: data}}
}
