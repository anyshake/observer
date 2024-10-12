package seisevent

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const EQZT_ID = "eqzt"

type EQZT struct {
	cache cache.AnyCache
}

func (j *EQZT) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      EQZT_ID,
		Country: "CN",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Zhao Tong Seismic Information System",
			"zh-TW": "昭通市地震資訊系統",
			"zh-CN": "昭通市地震信息系统",
		},
	}
}

func (j *EQZT) GetEvents(latitude, longitude float64) ([]Event, error) {
	if j.cache.Valid() {
		return j.cache.Get().([]Event), nil
	}

	res, err := request.GET(
		"http://www.eqzt.com/eqzt/seis/view_sbml.php?dzml=sbml",
		10*time.Second, time.Second, 3, false, nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	// Convert GB18030 to UTF-8
	dataBytes, err := io.ReadAll(transform.NewReader(bytes.NewReader(res), simplifiedchinese.GB18030.NewDecoder()))
	if err != nil {
		return nil, err
	}

	// Parse EQZT HTML response
	htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(dataBytes))
	if err != nil {
		return nil, err
	}

	var resultArr []Event
	htmlDoc.Find("table").Each(func(i int, s *goquery.Selection) {
		if val, ok := s.Attr("align"); !ok || val != "center" {
			return
		}
		s.Find("tr").Each(func(i int, s *goquery.Selection) {
			if i < 2 || i > 16 {
				return
			}

			var seisEvent Event
			s.Find("td").Each(func(i int, s *goquery.Selection) {
				if val, ok := s.Attr("colspan"); ok && val == "8" {
					return
				}
				textValue := s.Text()
				switch i {
				case 0:
					seisEvent.Depth = -1
					seisEvent.Verfied = true
					seisEvent.Event = fmt.Sprintf("No. %d", len(resultArr)+1)
					seisEvent.Timestamp = j.getTimestamp(textValue)
				case 1:
					seisEvent.Latitude = j.getLatitude(textValue)
				case 2:
					seisEvent.Longitude = j.getLongitude(textValue)
				case 3:
					seisEvent.Magnitude = append(seisEvent.Magnitude, j.getMagnitude("ML", textValue))
				case 4:
					seisEvent.Magnitude = append(seisEvent.Magnitude, j.getMagnitude("M", textValue))
				case 5:
					// Remove non-breaking space
					trimVal := strings.TrimFunc(textValue, func(r rune) bool { return r == ' ' })
					seisEvent.Region = trimVal
				}
			})
			seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
			seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

			resultArr = append(resultArr, seisEvent)
		})
	})

	sortedEvents := sortSeismicEvents(resultArr)
	j.cache.Set(sortedEvents)
	return sortedEvents, nil
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

func (j *EQZT) getMagnitude(magType, text string) Magnitude {
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	text = re.FindString(text)
	return Magnitude{
		Type:  ParseMagnitude(magType),
		Value: string2Float(text),
	}
}
