package seisevent

import (
	"bytes"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/corpix/uarand"
)

const BMKG_ID = "bmkg"

type BMKG struct {
	cache cache.AnyCache
}

func (c *BMKG) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      BMKG_ID,
		Country: "ID",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Meteorology, Climatology, and Geophysical Agency",
			"zh-TW": "印度尼西亞氣象、氣候和地球物理局",
			"zh-CN": "印度尼西亚气象、气候和地球物理局",
		},
	}
}

func (c *BMKG) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		return c.cache.Get().([]Event), nil
	}

	res, err := request.GET(
		"https://www.bmkg.go.id/gempabumi/gempabumi-dirasakan.bmkg",
		10*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	// Parse HTML response
	htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(res))
	if err != nil {
		return nil, err
	}

	var resultArr []Event
	htmlDoc.Find(".table-responsive").Each(func(i int, s *goquery.Selection) {
		s.Find("tbody").Each(func(i int, s *goquery.Selection) {
			s.Find("tr").Each(func(i int, s *goquery.Selection) {
				var seisEvent Event
				s.Find("td").Each(func(i int, s *goquery.Selection) {
					textValue := strings.TrimSpace(s.Text())
					switch i {
					case 0:
						seisEvent.Verfied = true
						seisEvent.Event = textValue
					case 1:
						seisEvent.Timestamp = c.getTimestamp(textValue)
					case 2:
						seisEvent.Latitude = c.getLatitude(textValue)
						seisEvent.Longitude = c.getLongitude(textValue)
					case 3:
						seisEvent.Magnitude = c.getMagnitude(textValue)
					case 4:
						seisEvent.Depth = c.getDepth(textValue)
					case 5:
						s.Find("a").Each(func(i int, s *goquery.Selection) {
							seisEvent.Region = strings.TrimSpace(s.Text())
						})
					}
				})
				seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
				seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

				resultArr = append(resultArr, seisEvent)
			})
		})
	})

	sortedEvents := sortSeismicEvents(resultArr)
	c.cache.Set(sortedEvents)
	return sortedEvents, nil
}

func (c *BMKG) getLatitude(data string) float64 {
	d := strings.Split(data, " ")
	if len(d) != 4 {
		return 0
	}

	return string2Float(d[0])
}

func (c *BMKG) getLongitude(data string) float64 {
	d := strings.Split(data, " ")
	if len(d) != 4 {
		return 0
	}

	return string2Float(d[2])
}

func (c *BMKG) getTimestamp(data string) int64 {
	t, _ := time.Parse("02/01/200615:04:05 WIB", data)
	return t.Add(-7 * time.Hour).UnixMilli()
}

func (c *BMKG) getDepth(data string) float64 {
	depthVal := strings.TrimSpace(strings.Replace(data, "Km", "", -1))
	return string2Float(depthVal)
}

func (c *BMKG) getMagnitude(data string) []Magnitude {
	return []Magnitude{
		{Type: ParseMagnitude("M"), Value: string2Float(data)},
	}
}
