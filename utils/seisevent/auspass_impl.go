package seisevent

import (
	"time"

	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const AUSPASS_ID = "auspass"

type AUSPASS struct {
	cache cache.AnyCache
}

func (c *AUSPASS) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      AUSPASS_ID,
		Country: "AU",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Australian Passive Seismic Server",
			"zh-TW": "澳洲國立大學地球科學研究學院",
			"zh-CN": "澳洲国立大学地球科学研究学院",
		},
	}
}

func (c *AUSPASS) GetEvents(latitude, longitude float64) ([]Event, error) {
	if !c.cache.Valid() {
		res, err := request.GET(
			"https://auspass.edu.au/fdsnws/event/1/query?minmag=-1&format=text&limit=300&orderby=time",
			60*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		c.cache.Set(res)
	}

	resultArr, err := ParseFdsnwsEvent(string(c.cache.Get().([]byte)), "2006-01-02T15:04:05", latitude, longitude)
	if err != nil {
		return nil, err
	}

	return sortSeismicEvents(resultArr), nil
}
