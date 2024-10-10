package seisevent

import (
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const BCSF_ID = "bcsf"

type BCSF struct {
	cache cache.AnyCache
}

func (c *BCSF) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      BCSF_ID,
		Country: "FR",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "French Central Seismological Office",
			"zh-TW": "法國中央地震局",
			"zh-CN": "法国中央地震局",
		},
	}
}

func (c *BCSF) GetEvents(latitude, longitude float64) ([]Event, error) {
	if !c.cache.Valid() {
		res, err := request.GET(
			"https://api.franceseisme.fr/fdsnws/event/1/query?minmag=-1&format=text&timezone=UTC&limit=300&orderby=time",
			30*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		c.cache.Set(res)
	}

	resultArr, err := ParseFdsnwsEvent(string(c.cache.Get().([]byte)), "2006-01-02T15:04:05.000000Z", latitude, longitude)
	if err != nil {
		return nil, err
	}

	return sortSeismicEvents(resultArr), nil
}
