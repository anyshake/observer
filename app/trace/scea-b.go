package trace

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/anyshake/observer/utils/duration"
	"github.com/anyshake/observer/utils/request"
)

type SCEA_B struct {
	DataSourceCache
}

func (s *SCEA_B) Property() string {
	return "四川地震局（速报）"
}

func (s *SCEA_B) Fetch() ([]byte, error) {
	if duration.Difference(time.Now(), s.Time) <= EXPIRATION {
		return s.Cache, nil
	}

	res, err := request.GET(
		"http://118.113.105.29:8002/api/bulletin/jsonPageList?pageSize=100",
		10*time.Second, time.Second, 3, false, nil,
	)
	if err != nil {
		return nil, err
	}

	s.Time = time.Now()
	s.Cache = make([]byte, len(res))
	copy(s.Cache, res)

	return res, nil
}

func (s *SCEA_B) Parse(data []byte) (map[string]any, error) {
	var result map[string]any
	err := json.Unmarshal(data, &result)
	if err != nil {
		fmt.Println(result)
		return nil, err
	}

	return result, nil
}

func (s *SCEA_B) Format(latitude, longitude float64, data map[string]any) ([]Event, error) {
	keys := []string{"eventId", "shockTime", "longitude", "latitude", "placeName", "magnitude", "depth"}

	var list []Event
	for _, v := range data["data"].([]any) {
		if !hasKey(v.(map[string]any), keys) || !isEmpty(v.(map[string]any), keys) {
			continue
		}

		l := Event{
			Verfied:   true,
			Depth:     v.(map[string]any)["depth"].(float64),
			Event:     v.(map[string]any)["eventId"].(string),
			Region:    v.(map[string]any)["placeName"].(string),
			Latitude:  v.(map[string]any)["latitude"].(float64),
			Longitude: v.(map[string]any)["longitude"].(float64),
			Magnitude: v.(map[string]any)["magnitude"].(float64),
			Timestamp: time.UnixMilli(int64(v.(map[string]any)["shockTime"].(float64))).UnixMilli(),
		}
		l.Distance = getDistance(latitude, l.Latitude, longitude, l.Longitude)
		l.Estimation = getEstimation(l.Depth, l.Distance)

		list = append(list, l)
	}

	return list, nil
}

func (s *SCEA_B) List(latitude, longitude float64) ([]Event, error) {
	res, err := s.Fetch()
	if err != nil {
		return nil, err
	}

	data, err := s.Parse(res)
	if err != nil {
		return nil, err
	}

	list, err := s.Format(latitude, longitude, data)
	if err != nil {
		return nil, err
	}

	return list, nil
}
