package trace

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type EQZT struct {
	dataSourceCache
}

func (j *EQZT) Property() string {
	return "昭通市地震信息系统"
}

func (j *EQZT) Fetch() ([]byte, error) {
	if time.Since(j.Time) <= EXPIRATION {
		return j.Cache, nil
	}

	res, err := request.GET(
		"http://www.eqzt.com/eqzt/seis/view_sbml.php?dzml=sbml",
		10*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	j.Time = time.Now()
	j.Cache = make([]byte, len(res))
	copy(j.Cache, res)

	return res, nil
}

func (j *EQZT) Parse(data []byte) (map[string]any, error) {
	result := make(map[string]any)
	result["data"] = make([]any, 0)

	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GB18030.NewDecoder())
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	htmlReader := bytes.NewBuffer(data)
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return nil, err
	}
	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		if val, ok := s.Attr("align"); !ok || val != "center" {
			return
		}
		s.Find("tr").Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				return
			}
			item := make(map[string]any)
			s.Find("td").Each(func(i int, s *goquery.Selection) {
				if val, ok := s.Attr("colspan"); ok && val == "8" {
					return
				}
				value := s.Text()
				switch i {
				case 0:
					item["timestamp"] = j.getTimestamp(value)
				case 1:
					item["latitude"] = j.getLatitude(value)
				case 2:
					item["longitude"] = j.getLongitude(value)
				case 4:
					item["magnitude"] = j.getMagnitude(value)
				case 5:
					trimVal := strings.TrimFunc(value, func(r rune) bool { return r == ' ' })
					item["event"] = trimVal
					item["region"] = fmt.Sprintf("云南及周边区域 - %s", trimVal)
				}
			})
			result["data"] = append(result["data"].([]any), item)
		})
	})

	return result, nil
}

func (j *EQZT) Format(latitude, longitude float64, data map[string]any) ([]seismicEvent, error) {
	var list []seismicEvent
	for _, v := range data["data"].([]any) {
		_, ok := v.(map[string]any)["timestamp"]
		if !ok {
			continue
		}
		l := seismicEvent{
			Verfied:   true,
			Depth:     -1,
			Latitude:  v.(map[string]any)["latitude"].(float64),
			Longitude: v.(map[string]any)["longitude"].(float64),
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

func (j *EQZT) List(latitude, longitude float64) ([]seismicEvent, error) {
	res, err := j.Fetch()
	if err != nil {
		return nil, err
	}

	data, err := j.Parse(res)
	if err != nil {
		return nil, err
	}

	list, err := j.Format(latitude, longitude, data)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (j *EQZT) getTimestamp(text string) int64 {
	t, _ := time.Parse("2006-01-02 15:04:05", text)
	return t.Add(-8 * time.Hour).UnixMilli()
}

func (j *EQZT) getLatitude(text string) float64 {
	text = strings.ReplaceAll(text, "°", " ")
	text = strings.ReplaceAll(text, "′", "")

	parts := strings.Fields(text)
	if len(parts) != 2 {
		return 0.0
	}

	degrees, err1 := strconv.ParseFloat(parts[0], 64)
	minutes, err2 := strconv.ParseFloat(parts[1], 64)
	if err1 != nil || err2 != nil {
		return 0.0
	}

	return degrees + minutes/60
}

func (j *EQZT) getLongitude(text string) float64 {
	text = strings.ReplaceAll(text, "°", " ")
	text = strings.ReplaceAll(text, "′", "")

	parts := strings.Fields(text)
	if len(parts) != 2 {
		return 0.0
	}

	degrees, err1 := strconv.ParseFloat(parts[0], 64)
	minutes, err2 := strconv.ParseFloat(parts[1], 64)
	if err1 != nil || err2 != nil {
		return 0.0
	}

	return degrees + minutes/60
}

func (j *EQZT) getMagnitude(text string) float64 {
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	text = re.FindString(text)
	return string2Float(text)
}
