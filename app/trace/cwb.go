package trace

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"com.geophone.observer/common/request"
	"github.com/PuerkitoBio/goquery"
)

type CWB struct{}

func (c *CWB) Property() (string, string) {
	const (
		NAME  string = "台湾交通部中央气象局"
		VALUE string = "CWB"
	)

	return NAME, VALUE
}

func (c *CWB) Fetch() ([]byte, error) {
	res, err := request.GETRequest(
		"https://www.cwb.gov.tw/V8/C/E/MOD/MAP_LIST.html",
		10*time.Second, time.Second, 3, false,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *CWB) Parse(data []byte) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	result["data"] = make([]interface{}, 0)

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
		item := make(map[string]interface{})

		item["latitude"] = latitude
		item["longitude"] = longitude
		item["depth"] = c.GetDepth(text)
		item["event"] = c.GetEvent(text)
		item["region"] = c.GetRegion(text)
		item["magnitude"] = c.GetMagnitude(text)
		item["timestamp"] = c.GetTimestamp(text)

		result["data"] = append(result["data"].([]interface{}), item)
	})

	return result, nil
}

func (c *CWB) Format(latitude, longitude float64, data map[string]interface{}) ([]EarthquakeList, error) {
	var list []EarthquakeList
	for _, v := range data["data"].([]interface{}) {
		l := EarthquakeList{
			Verfied:   true,
			Latitude:  String2Float(v.(map[string]interface{})["latitude"].(string)),
			Longitude: String2Float(v.(map[string]interface{})["longitude"].(string)),
			Depth:     v.(map[string]interface{})["depth"].(float64),
			Event:     v.(map[string]interface{})["event"].(string),
			Region:    v.(map[string]interface{})["region"].(string),
			Timestamp: v.(map[string]interface{})["timestamp"].(int64),
			Magnitude: v.(map[string]interface{})["magnitude"].(float64),
		}
		l.Distance = GetDistance(latitude, l.Latitude, longitude, l.Longitude)
		l.Estimated = GetEstimation(l.Distance)

		list = append(list, l)
	}

	return list, nil
}

func (c *CWB) List(latitude, longitude float64) ([]EarthquakeList, error) {
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

func (c *CWB) GetDepth(data string) float64 {
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

	return String2Float(result)
}

func (c *CWB) GetMagnitude(data string) float64 {
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
	return String2Float(result)
}

func (c *CWB) GetEvent(data string) string {
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

func (c *CWB) GetRegion(data string) string {
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

func (c *CWB) GetTimestamp(data string) int64 {
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

	return t.UnixMilli()
}
