package seisevent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
	"golang.org/x/sync/singleflight"
)

const ICL_ID = "icl"

type ICL struct {
	travelTimeTable *travel.AK135
	cache           cache.GenericCache[[]Event]
	sf              singleflight.Group
}

func (c *ICL) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      ICL_ID,
		Country: "CN",
		Default: "en-US",
		Locales: map[string]string{
			"en-US": "Institute of Care-life (ICL)",
			"zh-TW": "成都高新減災研究所",
		},
	}
}

func (h *ICL) GetEvents(ctx context.Context, latitude, longitude float64) ([]Event, error) {
	var baseEvents []Event

	if h.cache.Valid() {
		baseEvents = h.cache.Get()
	} else {
		v, err := fetchEvents(ctx, &h.sf, ICL_ID, func(ctx context.Context) ([]Event, error) {
			if h.cache.Valid() {
				return h.cache.Get(), nil
			}

			res, err := request.GET(
				ctx, "https://mobile-new.chinaeew.cn/v1/earlywarnings?start_at=&updates=",
				10*time.Second, time.Second, 3, false, nil,
				map[string]string{"User-Agent": uarand.GetRandom()},
			)
			if err != nil {
				return nil, err
			}

			// Parse ICL JSON response
			var dataMap map[string]any
			err = json.Unmarshal(res, &dataMap)
			if err != nil {
				return nil, err
			}

			// Check server response
			if dataMap["code"].(float64) != 0 {
				return nil, errors.New("remote server returned an error")
			}

			dataMapEvents, ok := dataMap["data"]
			if !ok {
				return nil, errors.New("seismic event data is not available")
			}

			// Ensure the response has the expected keys and they are not empty
			expectedKeys := []string{"eventId", "latitude", "longitude", "depth", "epicenter", "startAt", "magnitude"}

			var resultArr []Event
			for _, v := range dataMapEvents.([]any) {
				event := v.(map[string]any)

				if !isMapHasKeys(event, expectedKeys) || !isMapKeysEmpty(event, expectedKeys) {
					continue
				}

				resultArr = append(resultArr, Event{
					Verfied:   true,
					Depth:     event["depth"].(float64),
					Region:    event["epicenter"].(string),
					Latitude:  event["latitude"].(float64),
					Longitude: event["longitude"].(float64),
					Event:     fmt.Sprintf("%v", event["eventId"]),
					Magnitude: h.getMagnitude(event["magnitude"].(float64)),
					Timestamp: int64(event["startAt"].(float64)),
				})
			}

			sortedEvents := sortSeismicEvents(resultArr)
			h.cache.Set(sortedEvents)
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
			h.travelTimeTable,
			latitude,
			baseEvents[i].Latitude,
			longitude,
			baseEvents[i].Longitude,
			baseEvents[i].Depth,
		)
	}

	return baseEvents, nil
}

func (c *ICL) getMagnitude(data float64) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("M"), Value: data}}
}
