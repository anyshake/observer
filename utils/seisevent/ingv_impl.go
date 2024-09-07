package seisevent

import (
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const INGV_ID = "ingv"

type INGV struct {
	cache cache.BytesCache
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
	if !c.cache.Valid() {
		res, err := request.GET(
			"https://webservices.ingv.it/fdsnws/event/1/query?minmag=-1&format=text&timezone=UTC&limit=300&orderby=time",
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
