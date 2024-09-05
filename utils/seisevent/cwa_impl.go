package seisevent

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/anyshake/observer/utils/cache"
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

type CWA struct {
	cache cache.BytesCache
}

// Magic function that bypasses the Great Firewall of China
func (c *CWA) createGfwBypasser(customAddr string) *http.Transport {
	return &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, network, customAddr)
		},
	}
}

func (c *CWA) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      CWA_ID,
		Country: "TW",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Central Weather Administration (CWA)",
			"zh-TW": "交通部中央氣象署",
			"zh-CN": "交通部中央气象署",
		},
	}
}

func (c *CWA) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		res, err := request.GET(
			"https://www.cwa.gov.tw/V8/C/E/MOD/MAP_LIST.html",
			10*time.Second, time.Second, 3, false,
			c.createGfwBypasser("168.95.246.1:443"),
			map[string]string{"User-Agent": uarand.GetRandom()},
		)
		if err != nil {
			return nil, err
		}
		c.cache.Set(res)
	}

	// Parse CWA HTML response
	htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(c.cache.Get()))
	if err != nil {
		return nil, err
	}

	var resultArr []Event
	htmlDoc.Find("a").Each(func(i int, s *goquery.Selection) {
		eventLatitude, ok := s.Attr("data-lat")
		if !ok {
			return
		}
		eventLongitude, ok := s.Attr("data-lon")
		if !ok {
			return
		}

		textValue := s.Text()
		seisEvent := Event{
			Verfied:   true,
			Latitude:  string2Float(eventLatitude),
			Longitude: string2Float(eventLongitude),
			Depth:     c.getDepth(textValue),
			Event:     c.getEvent(textValue),
			Region:    c.getRegion(textValue),
			Magnitude: c.getMagnitude(textValue),
			Timestamp: c.getTimestamp(textValue),
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	})

	return sortSeismicEvents(resultArr), nil
}

func (c *CWA) getDepth(data string) float64 {
	exp := regexp.MustCompile(`深度(\d+(\.\d{1,}公里)|([1-9]\d*公里))`)
	if exp == nil {
		return -1
	}

	r := exp.FindAllStringSubmatch(data, -1)
	if len(r) == 0 || len(r[0]) == 0 {
		return -1
	}

	zh := regexp.MustCompile("[\u4e00-\u9fa5]+")
	result := zh.ReplaceAllString(r[0][0], "")

	return string2Float(result)
}

func (c *CWA) getMagnitude(data string) float64 {
	exp := regexp.MustCompile(`模\d+(\.\d{1,})|([1-9]\d*)$`)
	if exp == nil {
		return -1
	}

	r := exp.FindAllStringSubmatch(data, -1)
	if len(r) == 0 || len(r[0]) == 0 {
		return -1
	}

	zh := regexp.MustCompile("[\u4e00-\u9fa5]+")
	result := zh.ReplaceAllString(r[0][0], "")
	return string2Float(result)
}

func (c *CWA) getEvent(data string) string {
	exp := regexp.MustCompile(`地點為.+方\d+(\.\d{1,}公里)|([1-9]\d*公里)`)
	if exp == nil {
		return "未知地震"
	}

	r := exp.FindAllStringSubmatch(data, -1)
	if len(r) == 0 || len(r[0]) == 0 {
		return "未知地震"
	}

	result := strings.Replace(r[0][0], "地點為", "", -1)
	return result
}

func (c *CWA) getRegion(data string) string {
	exp := regexp.MustCompile(`\(位於.+\)`)
	if exp == nil {
		return "未知地点"
	}

	r := exp.FindAllStringSubmatch(data, -1)
	if len(r) == 0 || len(r[0]) == 0 {
		return "未知地点"
	}

	zh := regexp.MustCompile(`\(|\)|位於`)
	result := zh.ReplaceAllString(r[0][0], "")
	return result
}

func (c *CWA) getTimestamp(data string) int64 {
	exp := regexp.MustCompile(`時間為\d+月\d+日\d+時\d+，`)
	if exp == nil {
		return -1
	}

	r := exp.FindAllStringSubmatch(data, -1)
	if len(r) == 0 || len(r[0]) == 0 {
		return -1
	}

	zh := regexp.MustCompile("，|[\u4e00-\u9fa5]+")
	result := zh.ReplaceAllString(r[0][0], "")

	t, err := time.Parse("200601021504", fmt.Sprintf("%d%s", time.Now().Year(), result))
	if err != nil {
		return -1
	}

	return t.Add(-8 * time.Hour).UnixMilli()
}
