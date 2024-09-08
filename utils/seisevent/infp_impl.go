package seisevent

import (
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const INFP_ID = "infp"

type INFP struct {
	cache cache.BytesCache
}

func (c *INFP) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      INFP_ID,
		Country: "RO",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "National Institute for Earth Physics",
			"zh-TW": "國家地球物理研究所",
			"zh-CN": "国家地球物理研究所",
		},
	}
}

func (c *INFP) GetEvents(latitude, longitude float64) ([]Event, error) {
	if !c.cache.Valid() {
		res, err := request.GET(
			"https://eida-sc3.infp.ro/fdsnws/event/1/query?minmag=-1&format=text&limit=300&orderby=time",
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
