package seisevent

import (
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const SED_ID = "eida"

type SED struct {
	cache cache.BytesCache
}

func (c *SED) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      SED_ID,
		Country: "CH",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Swiss Seismological Service",
			"zh-TW": "瑞士地震局",
			"zh-CN": "瑞士地震局",
		},
	}
}

func (c *SED) GetEvents(latitude, longitude float64) ([]Event, error) {
	if !c.cache.Valid() {
		res, err := request.GET(
			"http://arclink.ethz.ch/fdsnws/event/1/query?minmag=-1&format=text&limit=300&orderby=time",
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
