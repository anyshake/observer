package trace

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bclswl0827/observer/utils/duration"
	"github.com/bclswl0827/observer/utils/request"
)

type USGS struct {
	DataSourceCache
}

func (u *USGS) Property() (string, string) {
	const (
		NAME  string = "United States Geological Survey"
		VALUE string = "USGS"
	)

	return NAME, VALUE
}

func (u *USGS) Fetch() ([]byte, error) {
	if duration.Difference(time.Now(), u.Time) <= EXPIRATION {
		return u.Cache, nil
	}

	res, err := request.GET(
		"https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/2.5_day.geojson",
		10*time.Second, time.Second, 3, false,
	)
	if err != nil {
		return nil, err
	}

	u.Time = time.Now()
	u.Cache = make([]byte, len(res))
	copy(u.Cache, res)

	return res, nil
}

func (u *USGS) Parse(data []byte) (map[string]any, error) {
	var result map[string]any
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *USGS) Format(latitude, longitude float64, data map[string]any) ([]Event, error) {
	events, ok := data["features"]
	if !ok {
		return nil, fmt.Errorf("source data is not valid")
	}

	var list []Event
	for _, v := range events.([]any) {
		if !hasKey(v.(map[string]any), []string{"properties"}) {
			continue
		}

		properties := v.(map[string]any)["properties"]
		if !hasKey(properties.(map[string]any), []string{
			"mag", "place", "time", "type", "title",
		}) {
			continue
		}

		geometry := v.(map[string]any)["geometry"]
		if !hasKey(geometry.(map[string]any), []string{"coordinates"}) {
			continue
		}

		coordinates := geometry.(map[string]any)["coordinates"]
		if len(coordinates.([]any)) != 3 {
			continue
		}

		if properties.(map[string]any)["type"].(string) != "earthquake" {
			continue
		}

		l := Event{
			Depth:     coordinates.([]any)[2].(float64),
			Verfied:   true,
			Timestamp: int64(properties.(map[string]any)["time"].(float64)),
			Event:     properties.(map[string]any)["title"].(string),
			Region:    properties.(map[string]any)["place"].(string),
			Latitude:  coordinates.([]any)[1].(float64),
			Longitude: coordinates.([]any)[0].(float64),
			Magnitude: properties.(map[string]any)["mag"].(float64),
		}
		l.Distance = getDistance(latitude, l.Latitude, longitude, l.Longitude)
		l.Estimation = getEstimation(l.Distance)

		list = append(list, l)
	}

	return list, nil
}

func (u *USGS) List(latitude, longitude float64) ([]Event, error) {
	res, err := u.Fetch()
	if err != nil {
		return nil, err
	}

	data, err := u.Parse(res)
	if err != nil {
		return nil, err
	}

	list, err := u.Format(latitude, longitude, data)
	if err != nil {
		return nil, err
	}

	return list, nil
}
