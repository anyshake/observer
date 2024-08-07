package trace

import (
	"encoding/json"
	"time"

	"github.com/anyshake/observer/utils/request"
)

type CEIC struct {
	dataSourceCache
}

func (c *CEIC) Property() string {
	return "中国地震台网中心"
}

func (c *CEIC) Fetch() ([]byte, error) {
	if time.Since(c.Time) <= EXPIRATION {
		return c.Cache, nil
	}

	res, err := request.GET(
		"https://news.ceic.ac.cn/ajax/google",
		10*time.Second, time.Second, 3, false, nil,
	)
	if err != nil {
		return nil, err
	}

	c.Time = time.Now()
	c.Cache = make([]byte, len(res))
	copy(c.Cache, res)

	return res, nil
}

func (c *CEIC) Parse(data []byte) (map[string]any, error) {
	result := make(map[string]any, 0)
	result["data"] = make([]any, 0)

	arr := make([]map[string]any, 0)
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return nil, err
	}

	result["data"] = arr
	return result, nil
}

func (c *CEIC) Format(latitude, longitude float64, data map[string]any) ([]seismicEvent, error) {
	keys := []string{"O_TIME", "EPI_LAT", "EPI_LON", "EPI_DEPTH", "M", "LOCATION_C"}

	var list []seismicEvent
	for _, v := range data["data"].([]map[string]any) {
		if !isMapHasKeys(v, keys) || !isMapKeysEmpty(v, keys) {
			continue
		}

		ts, err := time.Parse("2006-01-02 15:04:05", v["O_TIME"].(string))
		if err != nil {
			continue
		}

		// EPI_DEPTH type is not fixed
		var depth float64
		switch v["EPI_DEPTH"].(type) {
		case string:
			depth = string2Float(v["EPI_DEPTH"].(string))
		case float64:
			depth = v["EPI_DEPTH"].(float64)
		default:
			depth = -1
		}

		l := seismicEvent{
			Depth:     depth,
			Verfied:   true,
			Timestamp: ts.Add(-8 * time.Hour).UnixMilli(),
			Event:     v["LOCATION_C"].(string),
			Region:    v["LOCATION_C"].(string),
			Latitude:  string2Float(v["EPI_LAT"].(string)),
			Longitude: string2Float(v["EPI_LON"].(string)),
			Magnitude: string2Float(v["M"].(string)),
		}
		l.Distance = getDistance(latitude, l.Latitude, longitude, l.Longitude)
		l.Estimation = getSeismicEstimation(l.Depth, l.Distance)

		list = append(list, l)
	}

	return list, nil
}

func (c *CEIC) List(latitude, longitude float64) ([]seismicEvent, error) {
	res, err := c.Fetch()
	if err != nil {
		return nil, err
	}

	data, err := c.Parse(res)
	if err != nil {
		return nil, err
	}

	list, err := c.Format(latitude, longitude, data)
	if err != nil {
		return nil, err
	}

	return list, nil
}
