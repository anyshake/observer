package seisevent

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/corpix/uarand"
)

const SSN_ID = "ssn"

type SSN struct {
	cache cache.AnyCache
}

func (t *SSN) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      SSN_ID,
		Country: "MX",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "National Seismological Service",
			"zh-TW": "墨西哥國家地震局",
			"ja-JP": "メキシコ国立地震サービス局",
		},
	}
}

func (c *SSN) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		return c.cache.Get().([]Event), nil
	}

	res, err := request.GET(
		"http://www.ssn.unam.mx/sismicidad/ultimos-utc/",
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
	htmlDoc.Find(".table-latest").Each(func(i int, s *goquery.Selection) {
		s.Find("table").Each(func(i int, s *goquery.Selection) {
			s.Find("tr").Each(func(i int, s *goquery.Selection) {
				title, ok := s.Attr("title")
				if !ok {
					return
				}

				seisEvent := Event{Verfied: true, Event: title}
				s.Find("td").Each(func(i int, s *goquery.Selection) {
					switch i {
					case 0:
						mag := string2Float(s.Text())
						seisEvent.Magnitude = []Magnitude{{Type: ParseMagnitude("M"), Value: mag}}
					case 1:
						var dateStr, timeStr string
						s.Find("span").Each(func(i int, s *goquery.Selection) {
							switch i {
							case 0:
								dateStr = s.Text()
							case 1:
								timeStr = s.Text()
							}
						})
						ts, err := c.getTimestamp(fmt.Sprintf("%s %s", dateStr, timeStr))
						if err != nil {
							return
						}
						seisEvent.Timestamp = ts
					case 2:
						fields := strings.Split(s.Text(), ":")
						if len(fields) != 2 {
							return
						}
						seisEvent.Region = fields[0]
						coordinates := strings.Split(fields[1], ",")
						if len(coordinates) != 2 {
							return
						}
						seisEvent.Latitude, err = c.getLatitude(coordinates[0])
						if err != nil {
							return
						}
						seisEvent.Longitude, err = c.getLongitude(coordinates[1])
						if err != nil {
							return
						}
					case 3:
						depth, err := c.getDepth(s.Text())
						if err != nil {
							return
						}
						seisEvent.Depth = depth
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

func (t *SSN) getLatitude(latStr string) (float64, error) {
	latStr = strings.ReplaceAll(latStr, "°", "")
	latStr = strings.TrimSpace(latStr)
	return string2Float(latStr), nil
}

func (t *SSN) getLongitude(lonStr string) (float64, error) {
	lonStr = strings.ReplaceAll(lonStr, "°", "")
	lonStr = strings.TrimSpace(lonStr)
	return string2Float(lonStr), nil
}

func (t *SSN) getTimestamp(timeStr string) (int64, error) {
	tm, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return 0, err
	}
	return tm.UnixMilli(), nil
}

func (t *SSN) getDepth(depthStr string) (float64, error) {
	depthStr = strings.ReplaceAll(depthStr, "km", "")
	depthStr = strings.TrimSpace(depthStr)
	return string2Float(depthStr), nil
}
