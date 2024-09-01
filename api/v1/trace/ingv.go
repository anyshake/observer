package trace

import (
	"encoding/csv"
	"strconv"
	"strings"
	"time"

	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

type INGV struct {
	dataSourceCache
}

func (c *INGV) Property() string {
	return "Istituto nazionale di geofisica e vulcanologia"
}

func (c *INGV) Fetch() ([]byte, error) {
	if time.Since(c.Time) <= EXPIRATION {
		return c.Cache, nil
	}

	res, err := request.GET(
		"https://webservices.ingv.it/fdsnws/event/1/query?minmag=-1&format=text&timezone=UTC",
		10*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	c.Time = time.Now()
	c.Cache = make([]byte, len(res))
	copy(c.Cache, res)

	return res, nil
}

func (c *INGV) Parse(data []byte) (map[string]any, error) {
	result := make(map[string]any)
	result["data"] = make([]any, 0)

	csvDataStr := strings.ReplaceAll(string(data), "|", ",")
	reader := csv.NewReader(strings.NewReader(csvDataStr))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records[1:] {
		item := make(map[string]any)
		for i, v := range record {
			switch i {
			case 1:
				item["timestamp"] = c.getTimestamp(v)
			case 2:
				item["latitude"] = c.getLatitude(v)
			case 3:
				item["longitude"] = c.getLongitude(v)
			case 4:
				item["depth"] = c.getDepth(v)
			case 10:
				item["magnitude"] = c.getMagnitude(v)
			case 12:
				item["event"] = v
				item["region"] = v
			}
		}
		result["data"] = append(result["data"].([]any), item)
	}

	return result, nil
}

func (c *INGV) Format(latitude, longitude float64, data map[string]any) ([]seismicEvent, error) {
	var list []seismicEvent
	for _, v := range data["data"].([]any) {
		l := seismicEvent{
			Verfied:   true,
			Latitude:  v.(map[string]any)["latitude"].(float64),
			Longitude: v.(map[string]any)["longitude"].(float64),
			Depth:     v.(map[string]any)["depth"].(float64),
			Event:     v.(map[string]any)["event"].(string),
			Region:    v.(map[string]any)["region"].(string),
			Timestamp: v.(map[string]any)["timestamp"].(int64),
			Magnitude: v.(map[string]any)["magnitude"].(float64),
		}
		l.Distance = getDistance(latitude, l.Latitude, longitude, l.Longitude)
		l.Estimation = getSeismicEstimation(l.Depth, l.Distance)

		list = append(list, l)
	}

	return list, nil
}

func (c *INGV) List(latitude, longitude float64) ([]seismicEvent, error) {
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

func (c *INGV) getTimestamp(data string) int64 {
	t, _ := time.Parse("2006-01-02T15:04:05.000000", data)
	return t.UnixMilli()
}

func (c *INGV) getMagnitude(data string) float64 {
	m, err := strconv.ParseFloat(data, 64)
	if err == nil {
		return m
	}

	return 0
}

func (c *INGV) getDepth(data string) float64 {
	d, err := strconv.ParseFloat(data, 64)
	if err == nil {
		return d
	}

	return 0
}

func (c *INGV) getLatitude(data string) float64 {
	lat, err := strconv.ParseFloat(data, 64)
	if err == nil {
		return lat
	}

	return 0
}

func (c *INGV) getLongitude(data string) float64 {
	lng, err := strconv.ParseFloat(data, 64)
	if err == nil {
		return lng
	}

	return 0
}
