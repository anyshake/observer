package trace

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
	"github.com/anyshake/observer/utils/request"
	"github.com/corpix/uarand"
)

const HOST_IP_TO_BYPASS_GFW = "168.95.246.1:443"

type CWA struct {
	dataSourceCache
}

func (c *CWA) Property() string {
	return "交通部中央氣象署"
}

func (c *CWA) createGFWBypasser() *http.Transport {
	return &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, network, HOST_IP_TO_BYPASS_GFW)
		},
	}
}

func (c *CWA) Fetch() ([]byte, error) {
	if time.Since(c.Time) <= EXPIRATION {
		return c.Cache, nil
	}

	res, err := request.GET(
		"https://www.cwa.gov.tw/V8/C/E/MOD/MAP_LIST.html",
		10*time.Second, time.Second, 3, false,
		c.createGFWBypasser(),
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

func (c *CWA) Parse(data []byte) (map[string]any, error) {
	result := make(map[string]any)
	result["data"] = make([]any, 0)

	reader := bytes.NewBuffer(data)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		latitude, exists := s.Attr("data-lat")
		if !exists {
			return
		}

		longitude, exists := s.Attr("data-lon")
		if !exists {
			return
		}

		text := s.Text()
		item := make(map[string]any)

		item["latitude"] = latitude
		item["longitude"] = longitude
		item["depth"] = c.getDepth(text)
		item["event"] = c.getEvent(text)
		item["region"] = c.getRegion(text)
		item["magnitude"] = c.getMagnitude(text)
		item["timestamp"] = c.getTimestamp(text)

		result["data"] = append(result["data"].([]any), item)
	})

	return result, nil
}

func (c *CWA) Format(latitude, longitude float64, data map[string]any) ([]seismicEvent, error) {
	var list []seismicEvent
	for _, v := range data["data"].([]any) {
		l := seismicEvent{
			Verfied:   true,
			Latitude:  string2Float(v.(map[string]any)["latitude"].(string)),
			Longitude: string2Float(v.(map[string]any)["longitude"].(string)),
			Depth:     v.(map[string]any)["depth"].(float64),
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

func (c *CWA) List(latitude, longitude float64) ([]seismicEvent, error) {
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
