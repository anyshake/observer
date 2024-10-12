package seisevent

import (
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const USP_ID = "usp"

type USP struct {
	cache cache.AnyCache
}

func (c *USP) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      USP_ID,
		Country: "BR",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "USP Seismological Center",
			"zh-TW": "聖保羅大學地震學中心",
			"zh-CN": "圣保罗大学地震学中心",
		},
	}
}

func (c *USP) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		return c.cache.Get().([]Event), nil
	}

	res, err := request.GET(
		"https://moho.iag.usp.br/fdsnws/event/1/query?minmag=-1&format=text&limit=300&orderby=time",
		30*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	resultArr, err := ParseFdsnwsEvent(string(res), "2006-01-02T15:04:05", latitude, longitude)
	if err != nil {
		return nil, err
	}

	sortedEvents := sortSeismicEvents(resultArr)
	c.cache.Set(sortedEvents)
	return sortedEvents, nil
}
