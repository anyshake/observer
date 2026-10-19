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

const CENC_WOLFX_ID = "cenc_wolfx"

type CENC_WOLFX struct {
	travelTimeTable *travel.AK135
	cache           cache.GenericCache[[]Event]
	sf              singleflight.Group
}

func (c *CENC_WOLFX) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      CENC_WOLFX_ID,
		Country: "CN",
		Default: "en-US",
		Locales: map[string]string{
			"en-US": "China Earthquake Networks Center (Wolfx)",
			"zh-TW": "中國地震台網中心（Wolfx）",
		},
	}
}

func (c *CENC_WOLFX) GetEvents(ctx context.Context, latitude, longitude float64) ([]Event, error) {
	var baseEvents []Event

	if c.cache.Valid() {
		baseEvents = c.cache.Get()
	} else {
		v, err := fetchEvents(ctx, &c.sf, CENC_WOLFX_ID, func(ctx context.Context) ([]Event, error) {
			if c.cache.Valid() {
				return c.cache.Get(), nil
			}

			res, err := request.GET(
				ctx, "https://api.wolfx.jp/cenc_eqlist.json",
				10*time.Second, time.Second, 3, false, nil,
				map[string]string{"User-Agent": uarand.GetRandom()},
			)
			if err != nil {
				return nil, err
			}

			// Parse CENC JSON response
			var dataMapEvents map[string]any
			err = json.Unmarshal(res, &dataMapEvents)
			if err != nil {
				return nil, err
			}

			// Ensure the response has the expected keys and they are not empty
			expectedKeys := []string{"EventID", "type", "time", "location", "magnitude", "depth", "latitude", "longitude"}

			var resultArr []Event
			for _, v := range dataMapEvents {
				event, ok := v.(map[string]any)
				if !ok {
					continue
				}

				if !isMapHasKeys(event, expectedKeys) || !isMapKeysEmpty(event, expectedKeys) {
					continue
				}

				timestamp, err := c.getTimestamp(event["time"].(string))
				if err != nil {
					continue
				}

				resultArr = append(resultArr, Event{
					Verfied:   event["type"] == "reviewed",
					Event:     event["EventID"].(string),
					Timestamp: timestamp,
					Region:    event["location"].(string),
					Depth:     string2Float(event["depth"].(string)),
					Latitude:  string2Float(event["latitude"].(string)),
					Longitude: string2Float(event["longitude"].(string)),
					Magnitude: []Magnitude{{Type: "M", Value: string2Float(event["magnitude"].(string))}},
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

func (c *CENC_WOLFX) getTimestamp(timeStr string) (int64, error) {
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return 0, err
	}

	return t.Add(-8 * time.Hour).UnixMilli(), nil
}
