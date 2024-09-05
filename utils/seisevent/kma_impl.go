package seisevent

import (
	"bytes"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

type KMA struct {
	cache cache.BytesCache
}

func (k *KMA) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      KMA_ID,
		Country: "KR",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Korea Meteorological Administration",
			"zh-TW": "韓國氣象廳",
			"zh-CN": "韩国气象厅",
		},
	}
}

func (k *KMA) GetEvents(latitude, longitude float64) ([]Event, error) {
	if k.cache.Valid() {
		res, err := request.GET(
			"https://www.weather.go.kr/w/eqk-vol/search/korea.do",
			10*time.Second, time.Second, 3, false, nil,
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		k.cache.Set(res)
	}

	// Parse KMA HTML response
	htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(k.cache.Get()))
	if err != nil {
		return nil, err
	}

	var resultArr []Event
	htmlDoc.Find("#excel_body").Each(func(i int, s *goquery.Selection) {
		s.Find("tbody").Each(func(i int, s *goquery.Selection) {
			s.Find("tr").Each(func(i int, s *goquery.Selection) {
				var seisEvent Event

				s.Find("td").Each(func(i int, s *goquery.Selection) {
					textValue := strings.TrimSpace(s.Text())
					switch i {
					case 1:
						seisEvent.Verfied = true
						seisEvent.Timestamp = k.getTimestamp(textValue)
					case 2:
						seisEvent.Magnitude = k.getMagnitude(textValue)
					case 3:
						seisEvent.Depth = k.getDepth(textValue)
					case 5:
						seisEvent.Latitude = k.getLatitude(textValue)
					case 6:
						seisEvent.Longitude = k.getLongitude(textValue)
					case 7:
						seisEvent.Event = textValue
						seisEvent.Region = textValue
					}
				})
				seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
				seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

				resultArr = append(resultArr, seisEvent)
			})
		})
	})

	return resultArr, nil
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
