package seisevent

import (
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const KNMI_ID = "knmi"

type KNMI struct {
	cache cache.BytesCache
}

func (c *KNMI) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      KNMI_ID,
		Country: "NL",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Royal Netherlands Meteorological Institute",
			"zh-TW": "荷蘭皇家氣象研究所",
			"zh-CN": "荷兰皇家气象研究所",
		},
	}
}

func (c *KNMI) GetEvents(latitude, longitude float64) ([]Event, error) {
	if !c.cache.Valid() {
		res, err := request.GET(
			"https://rdsa.knmi.nl/fdsnws/event/1/query?minmag=-1&format=text&limit=300&orderby=time",
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
