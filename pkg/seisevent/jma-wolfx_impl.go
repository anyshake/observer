package seisevent

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
	"golang.org/x/sync/singleflight"
)

const JMA_WOLFX_ID = "jma_wolfx"

type JMA_WOLFX struct {
	travelTimeTable *travel.AK135
	cache           cache.GenericCache[[]Event]
	sf              singleflight.Group
}

func (j *JMA_WOLFX) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      JMA_WOLFX_ID,
		Country: "JP",
		Default: "en-US",
		Locales: map[string]string{
			"en-US": "Japan Meteorological Agency (Wolfx)",
			"zh-TW": "氣象廳（Wolfx）",
		},
	}
}

func (j *JMA_WOLFX) GetEvents(ctx context.Context, latitude, longitude float64) ([]Event, error) {
	var baseEvents []Event

	if j.cache.Valid() {
		baseEvents = j.cache.Get()
	} else {
		v, err := fetchEvents(ctx, &j.sf, JMA_WOLFX_ID, func(ctx context.Context) ([]Event, error) {
			if j.cache.Valid() {
				return j.cache.Get(), nil
			}

			res, err := request.GET(
				ctx, "https://api.wolfx.jp/jma_eqlist.json",
				10*time.Second, time.Second, 3, false, nil,
				map[string]string{"User-Agent": uarand.GetRandom()},
			)
			if err != nil {
				return nil, err
			}

			// Parse JMA_WOLFX JSON response
			var dataMapEvents map[string]any
			err = json.Unmarshal(res, &dataMapEvents)
			if err != nil {
				return nil, err
			}

			// Ensure the response has the expected keys and values
			expectedKeys := []string{"EventID", "time_full", "location", "magnitude", "depth", "latitude", "longitude"}

			var resultArr []Event
			for _, v := range dataMapEvents {
				event, ok := v.(map[string]any)
				if !ok {
					continue
				}

				if !isMapHasKeys(event, expectedKeys) || !isMapKeysEmpty(event, expectedKeys) {
					continue
				}

				timestamp, err := j.getTimestamp(event["time_full"].(string))
				if err != nil {
					continue
				}

				resultArr = append(resultArr, Event{
					Verfied:   true,
					Timestamp: timestamp,
					Depth:     j.getDepth(event["depth"].(string)),
					Event:     event["EventID"].(string),
					Region:    event["location"].(string),
					Latitude:  string2Float(event["latitude"].(string)),
					Longitude: string2Float(event["longitude"].(string)),
					Magnitude: j.getMagnitude(event["magnitude"].(string)),
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

func (j *JMA_WOLFX) getTimestamp(timeStr string) (int64, error) {
	t, err := time.Parse("2006/01/02 15:04:05", timeStr)
	if err != nil {
		return 0, err
	}

	return t.Add(-9 * time.Hour).UnixMilli(), nil
}

func (j *JMA_WOLFX) getDepth(data string) float64 {
	depthStr := strings.ReplaceAll(data, "km", "")
	return string2Float(depthStr)
}

func (j *JMA_WOLFX) getMagnitude(data string) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("M"), Value: string2Float(data)}}
}
