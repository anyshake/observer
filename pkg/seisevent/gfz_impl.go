package seisevent

import (
	"context"
	"slices"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
	"golang.org/x/sync/singleflight"
)

const GFZ_ID = "gfz"

type GFZ struct {
	travelTimeTable *travel.AK135
	cache           cache.GenericCache[[]Event]
	sf              singleflight.Group
}

func (c *GFZ) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      GFZ_ID,
		Country: "DE",
		Default: "en-US",
		Locales: map[string]string{
			"en-US": "GFZ German Research Centre",
			"zh-TW": "亥姆霍茲德國地理研究中心",
		},
	}
}

func (c *GFZ) GetEvents(ctx context.Context, latitude, longitude float64) ([]Event, error) {
	var baseEvents []Event

	if c.cache.Valid() {
		baseEvents = c.cache.Get()
	} else {
		v, err := fetchEvents(ctx, &c.sf, GFZ_ID, func(ctx context.Context) ([]Event, error) {
			if c.cache.Valid() {
				return c.cache.Get(), nil
			}

			res, err := request.GET(
				ctx, "https://geofon.gfz-potsdam.de/fdsnws/event/1/query?minmag=-1&format=text&limit=300&orderby=time",
				30*time.Second, time.Second, 3, false, nil,
				map[string]string{"User-Agent": uarand.GetRandom()},
			)
			if err != nil {
				return nil, err
			}

			resultArr, err := ParseFdsnwsEvent(string(res), "2006-01-02T15:04:05")
			if err != nil {
				return nil, err
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
