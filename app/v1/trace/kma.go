package trace

import (
	"bytes"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/anyshake/observer/utils/duration"
	"github.com/anyshake/observer/utils/request"
)

type KMA struct {
	DataSourceCache
}

func (k *KMA) Property() string {
	return "기상청（국내지진조회）"
}

func (k *KMA) Fetch() ([]byte, error) {
	if duration.Difference(time.Now(), k.Time) <= EXPIRATION {
		return k.Cache, nil
	}

	res, err := request.GET(
		"https://www.weather.go.kr/w/eqk-vol/search/korea.do",
		10*time.Second, time.Second, 3, false, nil,
	)
	if err != nil {
		return nil, err
	}

	k.Time = time.Now()
	k.Cache = make([]byte, len(res))
	copy(k.Cache, res)

	return res, nil
}

func (k *KMA) Parse(data []byte) (map[string]any, error) {
	result := make(map[string]any)
	result["data"] = make([]any, 0)

	reader := bytes.NewBuffer(data)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	doc.Find("#excel_body").Each(func(i int, s *goquery.Selection) {
		s.Find("tbody").Each(func(i int, s *goquery.Selection) {
			s.Find("tr").Each(func(i int, s *goquery.Selection) {
				item := make(map[string]any)
				s.Find("td").Each(func(i int, s *goquery.Selection) {
					value := strings.TrimSpace(s.Text())
					switch i {
					case 1:
						item["timestamp"] = k.getTimestamp(value)
					case 2:
						item["magnitude"] = k.getMagnitude(value)
					case 3:
						item["depth"] = k.getDepth(value)
					case 5:
						item["latitude"] = k.getLatitude(value)
					case 6:
						item["longitude"] = k.getLongitude(value)
					case 7:
						item["event"] = value
						item["region"] = value
					}
				})
				result["data"] = append(result["data"].([]any), item)
			})
		})
	})

	return result, nil
}

func (k *KMA) Format(latitude, longitude float64, data map[string]any) ([]Event, error) {
	var list []Event
	for _, v := range data["data"].([]any) {
		l := Event{
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
		l.Estimation = getEstimation(l.Depth, l.Distance)

		list = append(list, l)
	}

	return list, nil
}

func (k *KMA) List(latitude, longitude float64) ([]Event, error) {
	res, err := k.Fetch()
	if err != nil {
		return nil, err
	}

	data, err := k.Parse(res)
	if err != nil {
		return nil, err
	}

	list, err := k.Format(latitude, longitude, data)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (k *KMA) getTimestamp(data string) int64 {
	t, _ := time.Parse("2006/01/02 15:04:05", data)
	return t.Add(-9 * time.Hour).UnixMilli()
}

func (k *KMA) getMagnitude(data string) float64 {
	m, _ := strconv.ParseFloat(data, 64)
	return m
}

func (k *KMA) getDepth(data string) float64 {
	m, _ := strconv.ParseFloat(data, 64)
	return m
}

func (k *KMA) getLatitude(data string) float64 {
	numStr := strings.ReplaceAll(data, "N", "")
	numStr = strings.ReplaceAll(numStr, "S", "")
	numStr = strings.TrimSpace(numStr)

	if strings.Contains(data, "N") {
		longitude, _ := strconv.ParseFloat(numStr, 64)
		return longitude
	} else if strings.Contains(data, "S") {
		longitude, _ := strconv.ParseFloat(numStr, 64)
		return -longitude
	}

	return 0
}

func (k *KMA) getLongitude(data string) float64 {
	numStr := strings.ReplaceAll(data, "E", "")
	numStr = strings.ReplaceAll(numStr, "W", "")
	numStr = strings.TrimSpace(numStr)

	if strings.Contains(data, "E") {
		longitude, _ := strconv.ParseFloat(numStr, 64)
		return longitude
	} else if strings.Contains(data, "W") {
		longitude, _ := strconv.ParseFloat(numStr, 64)
		return -longitude
	}

	return 0
}
