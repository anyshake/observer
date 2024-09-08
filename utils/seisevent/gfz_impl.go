package seisevent

import (
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const GFZ_ID = "gfz"

type GFZ struct {
	cache cache.BytesCache
}

func (c *GFZ) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      GFZ_ID,
		Country: "DE",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "GFZ German Research Centre",
			"zh-TW": "亥姆霍茲德國地理研究中心",
			"zh-CN": "德国亥姆霍兹地球科学研究中心",
		},
	}
}

func (c *GFZ) GetEvents(latitude, longitude float64) ([]Event, error) {
	if !c.cache.Valid() {
		res, err := request.GET(
			"https://geofon.gfz-potsdam.de/fdsnws/event/1/query?minmag=-1&format=text&limit=300&orderby=time",
			30*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		c.cache.Set(res)
	}

	resultArr, err := ParseFdsnwsEvent(string(c.cache.Get()), "2006-01-02T15:04:05", latitude, longitude)
	if err != nil {
		return nil, err
	}

	return sortSeismicEvents(resultArr), nil
}
