package seisevent

import (
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/corpix/uarand"
)

const EMSC_ID = "emsc"

type EMSC struct {
	cache cache.AnyCache
}

func (c *EMSC) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      EMSC_ID,
		Country: "EU",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "European-Mediterranean Seismological Centre",
			"zh-TW": "歐洲與地中海地震中心",
			"zh-CN": "欧洲与地中海地震中心",
		},
	}
}

func (c *EMSC) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		return c.cache.Get().([]Event), nil
	}

	res, err := request.GET(
		"https://www.seismicportal.eu/fdsnws/event/1/query?minmag=-1&format=text&limit=300&orderby=time",
		10*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	resultArr, err := ParseFdsnwsEvent(string(res), "2006-01-02T15:04:05Z", latitude, longitude)
	if err != nil {
		return nil, err
	}

	sortedEvents := sortSeismicEvents(resultArr)
	c.cache.Set(sortedEvents)
	return sortedEvents, nil
}
