package trace

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

type SCEA_B struct {
	dataSourceCache
}

func (s *SCEA_B) Property() string {
	return "四川地震局（速报）"
}

func (s *SCEA_B) Fetch() ([]byte, error) {
	if time.Since(s.Time) <= EXPIRATION {
		return s.Cache, nil
	}

	res, err := request.GET(
		"http://118.113.105.29:8002/api/bulletin/jsonPageList?pageSize=100",
		10*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
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
		return nil, err
	}
	if result["code"].(float64) != 0 {
		return nil, fmt.Errorf("server error: %s", result["msg"])
	}

	return result, nil
}

func (s *SCEA_B) Format(latitude, longitude float64, data map[string]any) ([]seismicEvent, error) {
	keys := []string{"eventId", "shockTime", "longitude", "latitude", "placeName", "magnitude", "depth"}

	var list []seismicEvent
	for _, v := range data["data"].([]any) {
		if !isMapHasKeys(v.(map[string]any), keys) || !isMapKeysEmpty(v.(map[string]any), keys) {
			continue
		}

		l := seismicEvent{
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
		l.Estimation = getSeismicEstimation(l.Depth, l.Distance)

		list = append(list, l)
	}

	return list, nil
}

func (s *SCEA_B) List(latitude, longitude float64) ([]seismicEvent, error) {
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
