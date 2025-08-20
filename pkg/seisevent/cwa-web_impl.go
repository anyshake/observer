package seisevent

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/request"
	"github.com/corpix/uarand"
)

const CWA_WP_ID = "cwa_web"

type CWA_WP struct {
	cache cache.AnyCache
}

func (c *CWA_WP) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      CWA_WP_ID,
		Country: "TW",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Central Weather Administration (Web)",
			"zh-TW": "交通部中央氣象署（網頁）",
			"zh-CN": "交通部中央气象署（网页）",
		},
	}
}

func (c *CWA_WP) GetEvents(latitude, longitude float64) ([]Event, error) {
	if c.cache.Valid() {
		return c.cache.Get().([]Event), nil
	}

	// This is a workaround for the CWA webpage that does not provide the year of the events
	res, err := request.GET(
		"https://www.cloudflare.com/cdn-cgi/trace",
		10*time.Second, time.Second, 3, false,
		nil,
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}
	ts := regexp.MustCompile(`ts=(\d+)`).FindStringSubmatch(string(res))
	if len(ts) == 0 {
		return nil, errors.New("failed to get current time from Cloudflare")
	}
	currentYear := time.Unix(int64(string2Float(ts[1])), 0).Year()

	// Get CWA HTML response
	res, err = request.GET(
		"https://www.cwa.gov.tw/V8/C/E/MOD/MAP_LIST.html",
		10*time.Second, time.Second, 3, false,
		// Query CWA IP from custom encrypted DNS servers
		// Most overseas DoH / DoT providers are blocked in China
		// Recommended DNSCrypt: https://raw.githubusercontent.com/DNSCrypt/dnscrypt-resolvers/master/v3/public-resolvers.md
		(&CWA_API{}).createGfwBypasser([]string{
			// cs-tokyo
			"sdns://AQcAAAAAAAAADDE0Ni43MC4zMS40MyAxM3KtWVYywkFrhy8Jj4Ub3bllKExsvppPGQlkMNupWh4yLmRuc2NyeXB0LWNlcnQuY3J5cHRvc3Rvcm0uaXM",
			// dnscry.pt-tokyo-ipv4
			"sdns://AQcAAAAAAAAADDQ1LjY3Ljg2LjEyMyBDK5aRHZnKfdd6Q9ufEJY83WAQ9X5z7OAQa5CeptBCYBkyLmRuc2NyeXB0LWNlcnQuZG5zY3J5LnB0",
			// dnscry.pt-tokyo02-ipv4
			"sdns://AQcAAAAAAAAADDEwMy4xNzkuNDUuNiDfai5sp1im-BPHwbM1GCnTqn20FIbQfuvvybKsGf0pjhkyLmRuc2NyeXB0LWNlcnQuZG5zY3J5LnB0",
			// jp.tiar.app
			"sdns://AQcAAAAAAAAAEjE3Mi4xMDQuOTMuODA6MTQ0MyAyuHY-8b9lNqHeahPAzW9IoXnjiLaZpTeNbVs8TN9UUxsyLmRuc2NyeXB0LWNlcnQuanAudGlhci5hcHA",
			// saldns01-conoha-ipv4
			"sdns://gRQxNjMuNDQuMTI0LjIwNDo1MDQ0Mw",
			// saldns02-conoha-ipv4
			"sdns://gRUxNjAuMjUxLjIxNC4xNzI6NTA0NDM",
			// saldns03-conoha-ipv4
			"sdns://gRQxNjAuMjUxLjE2OC4yNTo1MDQ0Mw",
			// dnscry.pt-seoul-ipv4
			"sdns://AQcAAAAAAAAADTkyLjM4LjEzNS4xMjggyHfVGamJyxLfoAWjERmO4pY3KzKkqY-vSa2UnVx_gYAZMi5kbnNjcnlwdC1jZXJ0LmRuc2NyeS5wdA",
		}),
		map[string]string{"User-Agent": uarand.GetRandom()},
	)
	if err != nil {
		return nil, err
	}

	// Parse CWA HTML response
	htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(res))
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
		eventName, ok := s.Attr("data-name")
		if !ok {
			return
		}

		textValue := s.Text()
		seisEvent := Event{
			Verfied:   true,
			Event:     eventName,
			Latitude:  string2Float(eventLatitude),
			Longitude: string2Float(eventLongitude),
			Depth:     c.getDepth(textValue),
			Region:    c.getRegion(textValue),
			Magnitude: c.getMagnitude(textValue),
			Timestamp: c.getTimestamp(currentYear, textValue),
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	})

	sortedEvents := sortSeismicEvents(resultArr)
	c.cache.Set(sortedEvents)
	return sortedEvents, nil
}

func (c *CWA_WP) getDepth(data string) float64 {
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

func (c *CWA_WP) getMagnitude(data string) []Magnitude {
	exp := regexp.MustCompile(`模\d+(\.\d{1,})|([1-9]\d*)$`)
	if exp == nil {
		return []Magnitude{}
	}

	r := exp.FindAllStringSubmatch(data, -1)
	if len(r) == 0 || len(r[0]) == 0 {
		return []Magnitude{}
	}

	magnitude := string2Float(regexp.MustCompile("[\u4e00-\u9fa5]+").ReplaceAllString(r[0][0], ""))
	return []Magnitude{{Type: ParseMagnitude("ml"), Value: magnitude}}
}

func (c *CWA_WP) getRegion(data string) string {
	exp := regexp.MustCompile(`地點為.+方\d+(\.\d{1,}公里)|([1-9]\d*公里)`)
	if exp == nil {
		return "未知地點"
	}
	r := exp.FindAllStringSubmatch(data, -1)
	if len(r) == 0 || len(r[0]) == 0 {
		return "未知地點"
	}
	loc_1 := strings.Replace(r[0][0], "地點為", "", -1)

	exp = regexp.MustCompile(`\(位於.+\)`)
	if exp == nil {
		return loc_1
	}
	r = exp.FindAllStringSubmatch(data, -1)
	if len(r) == 0 || len(r[0]) == 0 {
		return loc_1
	}
	loc_2 := regexp.MustCompile(`\(|\)|位於`).ReplaceAllString(r[0][0], "")

	return fmt.Sprintf("%s (%s)", loc_1, loc_2)
}

func (c *CWA_WP) getTimestamp(year int, data string) int64 {
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

	t, err := time.Parse("200601021504", fmt.Sprintf("%d%s", year, result))
	if err != nil {
		return -1
	}

	return t.Add(-8 * time.Hour).UnixMilli()
}
