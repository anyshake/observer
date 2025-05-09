package seisevent

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/corpix/uarand"
)

const TMD_ID = "tmd"

type TMD struct {
	cache cache.AnyCache
}

func (t *TMD) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      TMD_ID,
		Country: "TH",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Thai Meteorological Department",
			"zh-TW": "泰國氣象局",
			"ja-JP": "タイ気象局",
		},
	}
}

func (c *TMD) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		return c.cache.Get().([]Event), nil
	}

	res, err := request.GET(
		"https://earthquake.tmd.go.th/inside.html?ps=200",
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
	htmlDoc.Find(".tbis_leq1").Each(func(eventIdx int, s *goquery.Selection) {
		seisEvent := Event{Verfied: true}
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			switch i {
			case 0:
				timeStr := s.Text()
				if len(timeStr) >= 19 {
					timeStr = timeStr[:19]
				}
				ts, err := c.getTimestamp(timeStr)
				if err != nil {
					return
				}
				seisEvent.Timestamp = ts
			case 1:
				mag := string2Float(s.Text())
				seisEvent.Magnitude = []Magnitude{{Type: ParseMagnitude("M"), Value: mag}}
			case 2:
				seisEvent.Latitude = string2Float(s.Text())
			case 3:
				seisEvent.Longitude = string2Float(s.Text())
			case 4:
				seisEvent.Depth = string2Float(s.Text())
			case 6:
				thai, eng, err := getRegion(s)
				if err != nil {
					return
				}
				seisEvent.Event = strconv.Itoa(eventIdx)
				seisEvent.Region = fmt.Sprintf("%s (%s)", thai, eng)
			}
		})

		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	})
	htmlDoc.Find(".tbis_leq2").Each(func(eventIdx int, s *goquery.Selection) {
		seisEvent := Event{Verfied: true}
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			switch i {
			case 0:
				timeStr := s.Text()
				if len(timeStr) >= 19 {
					timeStr = timeStr[:19]
				}
				ts, err := c.getTimestamp(timeStr)
				if err != nil {
					return
				}
				seisEvent.Timestamp = ts
			case 1:
				mag := string2Float(s.Text())
				seisEvent.Magnitude = []Magnitude{{Type: ParseMagnitude("M"), Value: mag}}
			case 2:
				seisEvent.Latitude = string2Float(s.Text())
			case 3:
				seisEvent.Longitude = string2Float(s.Text())
			case 4:
				seisEvent.Depth = string2Float(s.Text())
			case 6:
				thai, eng, err := getRegion(s)
				if err != nil {
					return
				}
				seisEvent.Event = strconv.Itoa(eventIdx)
				seisEvent.Region = fmt.Sprintf("%s (%s)", thai, eng)
			}
		})

		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	})

	sortedEvents := sortSeismicEvents(resultArr)
	c.cache.Set(sortedEvents)
	return sortedEvents, nil
}

func getRegion(s *goquery.Selection) (string, string, error) {
	var thai, eng string
	s.Contents().Each(func(i int, s *goquery.Selection) {
		html, _ := s.Html()
		region := strings.Split(html, "<br/>")
		if len(region) == 2 {
			thai = region[0]
			eng = region[1]
		}
	})

	if thai == "" || eng == "" {
		return "", "", errors.New("failed to parse region")
	}

	return thai, eng, nil
}

func (t *TMD) getTimestamp(timeStr string) (int64, error) {
	tm, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return 0, err
	}
	tm = tm.Add(-7 * time.Hour)
	return tm.UnixMilli(), nil
}
