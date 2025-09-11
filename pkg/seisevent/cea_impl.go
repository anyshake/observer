package seisevent

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/bclswl0827/travel"
	"github.com/corpix/uarand"
)

const CEA_ID = "cea"

type CEA struct {
	travelTimeTable *travel.AK135
	cache           cache.AnyCache
}

func (c *CEA) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      CEA_ID,
		Country: "FR",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "CEA/Dase - Earth and Environmental Science",
			"zh-TW": "原子能和替代能源委員會",
			"zh-CN": "原子能和替代能源委员会",
		},
	}
}

func (c *CEA) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		return c.cache.Get().([]Event), nil
	}

	res, err := request.GET(
		"https://www-dase.cea.fr/evenement/derniers_evenements.php",
		10*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	// Parse CEA/DASE HTML response
	htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(res))
	if err != nil {
		return nil, err
	}

	var resultArr []Event
	htmlDoc.Find(".arial11bleu").Each(func(i int, s *goquery.Selection) {
		var (
			dateText  string
			seisEvent Event
		)
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			textValue := strings.TrimSpace(s.Text())
			switch i {
			case 0:
				seisEvent.Verfied = false
				seisEvent.Depth = -1
				dateText = textValue
			case 1:
				fullTimeString := fmt.Sprintf("%s %s", dateText, textValue)
				seisEvent.Timestamp = c.getTimestamp(fullTimeString)
			case 2:
				seisEvent.Latitude = c.getLatitude(textValue)
				seisEvent.Longitude = c.getLongitude(textValue)
			case 3:
				seisEvent.Event = textValue
				seisEvent.Region = textValue
			case 4:
				seisEvent.Magnitude = c.getMagnitude(textValue)
			}
		})
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(c.travelTimeTable, latitude, seisEvent.Latitude, longitude, seisEvent.Longitude, seisEvent.Depth)

		resultArr = append(resultArr, seisEvent)
	})

	sortedEvents := sortSeismicEvents(resultArr)
	c.cache.Set(sortedEvents)
	return sortedEvents, nil
}

func (c *CEA) getTimestamp(data string) int64 {
	t, _ := time.Parse("02/01/2006 15:04:05", data)
	return t.Add(-2 * time.Hour).UnixMilli()
}

func (c *CEA) getMagnitude(data string) []Magnitude {
	m := strings.Split(data, "=")
	if len(m) > 1 {
		magnitudeType := strings.ToUpper(m[0])
		magnitudeVal, _ := strconv.ParseFloat(m[1], 64)
		return []Magnitude{
			{Type: ParseMagnitude(magnitudeType), Value: magnitudeVal},
		}
	}

	return []Magnitude{}
}

func (c *CEA) getLatitude(data string) float64 {
	pos := strings.Split(data, ",")
	if len(pos) > 1 {
		latitude, _ := strconv.ParseFloat(pos[0], 64)
		return latitude
	}

	return 0
}

func (c *CEA) getLongitude(data string) float64 {
	pos := strings.Split(data, ",")
	if len(pos) > 1 {
		longitude, _ := strconv.ParseFloat(pos[1], 64)
		return longitude
	}

	return 0
}
