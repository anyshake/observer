package fdsnws

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/testsupport"
	fdsnwsgo "github.com/bclswl0827/fdsnws-go"
)

func TestDataSelectQuality(t *testing.T) {
	db, handler := testsupport.OpenDAO(t)
	testsupport.InitStation(t, handler)
	start := time.Date(2026, 10, 19, 12, 0, 0, 0, time.UTC)
	end := start.Add(time.Second)
	testsupport.EnsureSeisTable(t, db, start)
	dataProvider := &provider{handler: handler, device: testsupport.NewHardware(), now: func() time.Time { return start }}
	server, err := fdsnwsgo.NewFDSNWS(dataProvider,
		fdsnwsgo.WithModules(fdsnwsgo.ModuleDataSelect),
		fdsnwsgo.WithService(newDataSelectService(dataProvider)),
	)
	if err != nil {
		t.Fatal(err)
	}
	for _, method := range []string{http.MethodGet, http.MethodPost} {
		for _, quality := range []string{"", "B", "b", "D", "d", "R", "Q", "M", "Z"} {
			name := quality
			if name == "" {
				name = "default"
			}
			t.Run(method+"/"+name, func(t *testing.T) {
				request := fdsnwsgo.Request{Method: method, Path: "/fdsnws/dataselect/1/query"}
				if method == http.MethodGet {
					values := url.Values{
						"network": {"*"}, "station": {"*"}, "location": {"*"}, "channel": {"*"},
						"starttime": {start.Format(time.RFC3339)}, "endtime": {end.Format(time.RFC3339)},
					}
					if quality != "" {
						values.Set("QUALITY", quality)
					}
					request.RawQuery = values.Encode()
				} else {
					body := "* * * * " + start.Format(time.RFC3339) + " " + end.Format(time.RFC3339)
					if quality != "" {
						body = "QUALITY = " + quality + "\n" + body
					}
					request.Body = []byte(body)
				}
				response := server.Handle(context.Background(), request)
				wantStatus := http.StatusNoContent
				if quality == "Z" {
					wantStatus = http.StatusBadRequest
				}
				if response.StatusCode != wantStatus {
					t.Fatalf("status = %d; want %d: %s", response.StatusCode, wantStatus, response.Body)
				}
				if wantStatus == http.StatusBadRequest && !strings.HasPrefix(string(response.Body), "Error 400: Bad request\n") {
					t.Fatalf("quality error did not use the FDSN error format: %s", response.Body)
				}
			})
		}
	}
}
