package trace

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

type CEA_DASE struct {
	dataSourceCache
}

func (c *CEA_DASE) Property() string {
	return "Commissariat à l'Energie Atomique"
}

func (c *CEA_DASE) Fetch() ([]byte, error) {
	if time.Since(c.Time) <= EXPIRATION {
		return c.Cache, nil
	}

	res, err := request.GET(
		"https://www-dase.cea.fr/evenement/derniers_evenements.php",
		10*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	c.Time = time.Now()
	c.Cache = make([]byte, len(res))
	copy(c.Cache, res)

	return res, nil
}

func (c *CEA_DASE) Parse(data []byte) (map[string]any, error) {
	result := make(map[string]any)
	result["data"] = make([]any, 0)

	reader := bytes.NewBuffer(data)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	doc.Find(".arial11bleu").Each(func(i int, s *goquery.Selection) {
		item := make(map[string]any)
		var dateString string
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			value := strings.TrimSpace(s.Text())
			switch i {
			case 0:
				item["depth"] = -1
				dateString = value
			case 1:
				item["timestamp"] = c.getTimestamp(fmt.Sprintf("%s %s", dateString, value))
			case 2:
				item["latitude"] = c.getLatitude(value)
				item["longitude"] = c.getLongitude(value)
			case 3:
				item["event"] = value
				item["region"] = value
			case 4:
				item["magnitude"] = c.getMagnitude(value)
			}
		})
		result["data"] = append(result["data"].([]any), item)
	})

	return result, nil
}

func (c *CEA_DASE) Format(latitude, longitude float64, data map[string]any) ([]seismicEvent, error) {
	var list []seismicEvent
	for _, v := range data["data"].([]any) {
		l := seismicEvent{
			Verfied:   false,
			Latitude:  v.(map[string]any)["latitude"].(float64),
			Longitude: v.(map[string]any)["longitude"].(float64),
			Depth:     float64(v.(map[string]any)["depth"].(int)),
			Event:     v.(map[string]any)["event"].(string),
			Region:    v.(map[string]any)["region"].(string),
			Timestamp: v.(map[string]any)["timestamp"].(int64),
			Magnitude: v.(map[string]any)["magnitude"].(float64),
		}
		l.Distance = getDistance(latitude, l.Latitude, longitude, l.Longitude)
		l.Estimation = getSeismicEstimation(l.Depth, l.Distance)

		list = append(list, l)
	}

	return list, nil
}

func (c *CEA_DASE) List(latitude, longitude float64) ([]seismicEvent, error) {
	res, err := c.Fetch()
	if err != nil {
		return nil, err
	}

	data, err := c.Parse(res)
	if err != nil {
		return nil, err
	}

	list, err := c.Format(latitude, longitude, data)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (c *CEA_DASE) getTimestamp(data string) int64 {
	t, _ := time.Parse("02/01/2006 15:04:05", data)
	return t.Add(-1 * time.Hour).UnixMilli()
}

func (c *CEA_DASE) getMagnitude(data string) float64 {
	m := strings.Split(data, "=")
	if len(m) > 1 {
		magnitude, _ := strconv.ParseFloat(m[1], 64)
		return magnitude
	}

	return 0
}

func (c *CEA_DASE) getLatitude(data string) float64 {
	pos := strings.Split(data, ",")
	if len(pos) > 1 {
		latitude, _ := strconv.ParseFloat(pos[0], 64)
		return latitude
	}

	return 0
}

func (c *CEA_DASE) getLongitude(data string) float64 {
	pos := strings.Split(data, ",")
	if len(pos) > 1 {
		longitude, _ := strconv.ParseFloat(pos[1], 64)
		return longitude
	}

	return 0
}
