package seisevent

import (
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
)

const INGV_ID = "ingv"

type INGV struct {
	travelTimeTable *travel.AK135
	cache           cache.AnyCache
}

func (c *INGV) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      INGV_ID,
		Country: "IT",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "National Institute of Geophysics and Volcanology",
			"zh-TW": "國立地球物理與火山學研究所",
			"zh-CN": "国立地球物理与火山学研究所",
		},
	}
}

func (c *INGV) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		return c.cache.Get().([]Event), nil
	}

	res, err := request.GET(
		"https://webservices.ingv.it/fdsnws/event/1/query?minmag=-1&format=text&timezone=UTC&limit=300&orderby=time",
		30*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	resultArr, err := ParseFdsnwsEvent(c.travelTimeTable, string(res), "2006-01-02T15:04:05", latitude, longitude)
	if err != nil {
		return nil, err
	}

	sortedEvents := sortSeismicEvents(resultArr)
	c.cache.Set(sortedEvents)
	return sortedEvents, nil
}
