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

const GEONET_ID = "geonet"

type GEONET struct {
	travelTimeTable *travel.AK135
	cache           cache.GenericCache[[]Event]
	sf              singleflight.Group
}

func (u *GEONET) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      GEONET_ID,
		Country: "NZ",
		Default: "en-US",
		Locales: map[string]string{
			"en-US": "The GeoNet Project",
			"zh-TW": "GeoNet 計畫",
		},
	}
}

func (c *GEONET) GetEvents(ctx context.Context, latitude, longitude float64) ([]Event, error) {
	var baseEvents []Event

	if c.cache.Valid() {
		baseEvents = c.cache.Get()
	} else {
		v, err := fetchEvents(ctx, &c.sf, GEONET_ID, func(ctx context.Context) ([]Event, error) {
			if c.cache.Valid() {
				return c.cache.Get(), nil
			}

			res, err := request.GET(
				ctx, "https://api.geonet.org.nz/quake?MMI=1",
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

				resultArr = append(resultArr, Event{
					Verfied:   true,
					Latitude:  coordinates.([]any)[1].(float64),
					Longitude: coordinates.([]any)[0].(float64),
					Depth:     properties.(map[string]any)["depth"].(float64),
					Event:     properties.(map[string]any)["publicID"].(string),
					Region:    properties.(map[string]any)["locality"].(string),
					Magnitude: c.getMagnitude(properties.(map[string]any)["magnitude"].(float64)),
					Timestamp: c.getTimestamp(properties.(map[string]any)["time"].(string)),
				})
			}

			sortedEvents := sortSeismicEvents(resultArr)
			c.cache.Set(sortedEvents)
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
			c.travelTimeTable,
			latitude,
			baseEvents[i].Latitude,
			longitude,
			baseEvents[i].Longitude,
			baseEvents[i].Depth,
		)
	}

	return baseEvents, nil
}

func (u *GEONET) getTimestamp(data string) int64 {
	t, _ := time.Parse("2006-01-02T15:04:05.000Z", data)
	return t.UnixMilli()
}

func (u *GEONET) getMagnitude(data float64) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("M"), Value: data}}
}
