package trace

import (
	"time"

	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

type SCEA_E struct {
	SCEA_B
	dataSourceCache
}

func (s *SCEA_E) Property() string {
	return "四川地震局（预警）"
}

func (s *SCEA_E) Fetch() ([]byte, error) {
	if time.Since(s.Time) <= EXPIRATION {
		return s.Cache, nil
	}

	res, err := request.GET(
		"http://118.113.105.29:8002/api/earlywarning/jsonPageList?pageSize=100",
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

func (s *SCEA_E) Format(latitude, longitude float64, data map[string]any) ([]seismicEvent, error) {
	keys := []string{"eventId", "shockTime", "longitude", "latitude", "placeName", "magnitude", "depth"}

	var list []seismicEvent
	for _, v := range data["data"].([]any) {
		if !isMapHasKeys(v.(map[string]any), keys) || !isMapKeysEmpty(v.(map[string]any), keys) {
			continue
		}

		l := seismicEvent{
			Verfied:   true,
			Depth:     -1,
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

func (s *SCEA_E) List(latitude, longitude float64) ([]seismicEvent, error) {
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
