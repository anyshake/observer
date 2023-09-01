package trace

import (
	"time"

	"com.geophone.observer/utils/request"
	t "com.geophone.observer/utils/time"
)

type SCEA_E struct {
	SCEA_B
	DataSourceCache
}

func (s *SCEA_E) Property() (string, string) {
	const (
		NAME  string = "四川地震局预警"
		VALUE string = "SCEA-E"
	)

	return NAME, VALUE
}

func (s *SCEA_E) Fetch() ([]byte, error) {
	if t.Diff(time.Now(), s.Time) <= EXPIRATION {
		return s.Cache, nil
	}

	res, err := request.GET(
		"http://118.113.105.29:8002/api/earlywarning/jsonPageList?pageSize=100",
		10*time.Second, time.Second, 3, false,
	)
	if err != nil {
		return nil, err
	}

	s.Time = time.Now()
	s.Cache = make([]byte, len(res))
	copy(s.Cache, res)

	return res, nil
}
