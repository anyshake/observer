package seisevent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/dnsquery"
	"github.com/anyshake/observer/pkg/request"
	"github.com/corpix/uarand"
	"github.com/miekg/dns"
	"golang.org/x/exp/rand"
)

const CWA_API_ID = "cwa_api"

type CWA_API struct {
	cache cache.AnyCache
}

// Magic function that bypasses the Great Firewall of China
func (c *CWA_API) createGfwBypasser(dnsList []string) *http.Transport {
	return &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			hostname, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse address: %w", err)
			}

			dnsResolver, err := dnsquery.New(dnsList[rand.Intn(len(dnsList))])
			if err != nil {
				return nil, fmt.Errorf("failed to create DNS resolver: %w", err)
			}

			if err := dnsResolver.Open(); err != nil {
				return nil, fmt.Errorf("failed to open DNS resolver: %w", err)
			}
			defer dnsResolver.Close()

			res, err := dnsResolver.Query((&dns.Msg{}).SetQuestion(fmt.Sprintf("%s.", hostname), dns.TypeA), 5*time.Second)
			if err != nil {
				return nil, fmt.Errorf("failed to query DNS: %w", err)
			}
			if len(res.Answer) == 0 {
				return nil, errors.New("no answer from DNS")
			}

			dnsAnswer := res.Answer[rand.Intn(len(res.Answer))]
			txtRecord, ok := dnsAnswer.(*dns.A)
			if !ok {
				return nil, fmt.Errorf("unexpected DNS answer type: %T", dnsAnswer)
			}
			if len(txtRecord.A) == 0 {
				return nil, errors.New("no answer from DNS")
			}

			return (&net.Dialer{}).DialContext(ctx, network, fmt.Sprintf("%s:%s", txtRecord.A.String(), port))
		},
	}
}

func (c *CWA_API) getRequestBody(limit int) string {
	return fmt.Sprintf(`draw=1&columns[0][data]=0&columns[0][name]=EventNo&columns[0][searchable]=false&columns[0][orderable]=true&columns[0][search][value]=&columns[0][search][regex]=false&columns[1][data]=1&columns[1][name]=MaxIntensity&columns[1][searchable]=true&columns[1][orderable]=true&columns[1][search][value]=&columns[1][search][regex]=false&columns[2][data]=2&columns[2][name]=OriginTime&columns[2][searchable]=true&columns[2][orderable]=true&columns[2][search][value]=&columns[2][search][regex]=false&columns[3][data]=3&columns[3][name]=MagnitudeValue&columns[3][searchable]=true&columns[3][orderable]=true&columns[3][search][value]=&columns[3][search][regex]=false&columns[4][data]=4&columns[4][name]=Depth&columns[4][searchable]=true&columns[4][orderable]=true&columns[4][search][value]=&columns[4][search][regex]=false&columns[5][data]=5&columns[5][name]=Description&columns[5][searchable]=true&columns[5][orderable]=true&columns[5][search][value]=&columns[5][search][regex]=false&columns[6][data]=6&columns[6][name]=Description&columns[6][searchable]=true&columns[6][orderable]=true&columns[6][search][value]=&columns[6][search][regex]=false&order[0][column]=2&order[0][dir]=desc&start=0&length=%d&search[value]=&search[regex]=false&Search=&txtSDate=&txtEDate=&txtSscale=&txtEscale=&txtSdepth=&txtEdepth=&txtLonS=&txtLonE=&txtLatS=&txtLatE=&ddlCity=&ddlCitySta=&txtIntensityB=&txtIntensityE=&txtLon=&txtLat=&txtKM=&ddlStationName=&cblEventNo=&txtSDatePWS=&txtEDatePWS=&txtSscalePWS=&txtEscalePWS=&ddlMark=`, limit)
}

func (c *CWA_API) GetProperty() DataSourceProperty {
	return DataSourceProperty{
		ID:      CWA_API_ID,
		Country: "TW",
		Deafult: "en-US",
		Locales: map[string]string{
			"en-US": "Central Weather Administration (API)",
			"zh-TW": "交通部中央氣象署（API）",
			"zh-CN": "交通部中央气象署（API）",
		},
	}
}

func (c *CWA_API) GetEvents(latitude, longitude float64) ([]Event, error) {
	// Get CWA HTML response
	if c.cache.Valid() {
		return c.cache.Get().([]Event), nil
	}

	res, err := request.POST(
		"https://scweb.cwa.gov.tw/zh-tw/earthquake/ajaxhandler",
		c.getRequestBody(100),
		"application/x-www-form-urlencoded; charset=UTF-8",
		10*time.Second, time.Second, 3, false,
		// Query CWA IP from custom encrypted DNS servers
		// Most overseas DoH / DoT providers are blocked in China
		// Recommended DNSCrypt: https://raw.githubusercontent.com/DNSCrypt/dnscrypt-resolvers/master/v3/public-resolvers.md
		c.createGfwBypasser([]string{
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

	// Parse CWA JSON response
	var dataMap map[string]any
	err = json.Unmarshal(res, &dataMap)
	if err != nil {
		return nil, err
	}

	dataMapEvents, ok := dataMap["data"].([]any)
	if !ok {
		return nil, errors.New("seismic event data is not available")
	}

	var resultArr []Event
	for _, event := range dataMapEvents {
		eventData := event.([]any)
		if len(eventData) < 10 {
			continue
		}

		seisEvent := Event{
			Verfied:   true,
			Depth:     string2Float(eventData[4].(string)),
			Event:     eventData[0].(string),
			Region:    eventData[5].(string),
			Latitude:  string2Float(eventData[8].(string)),
			Longitude: string2Float(eventData[7].(string)),
			Magnitude: c.getMagnitude(eventData[3].(string)),
			Timestamp: c.getTimestamp(eventData[2].(string)),
		}
		seisEvent.Distance = getDistance(latitude, seisEvent.Latitude, longitude, seisEvent.Longitude)
		seisEvent.Estimation = getSeismicEstimation(seisEvent.Depth, seisEvent.Distance)

		resultArr = append(resultArr, seisEvent)
	}

	sortedEvents := sortSeismicEvents(resultArr)
	c.cache.Set(sortedEvents)
	return sortedEvents, nil
}

func (c *CWA_API) getTimestamp(textValue string) int64 {
	t, _ := time.Parse("2006-01-02 15:04:05", textValue)
	return t.Add(-8 * time.Hour).UnixMilli()
}

func (c *CWA_API) getMagnitude(textValue string) []Magnitude {
	return []Magnitude{{Type: ParseMagnitude("ML"), Value: string2Float(textValue)}}
}
