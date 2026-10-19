package fdsnws

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/testsupport"
	fdsnwsgo "github.com/bclswl0827/fdsnws-go"
	"github.com/bclswl0827/fdsnws-go/service/station"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestSetupAccessKeyAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, handler := testsupport.OpenDAO(t)
	if err := db.AutoMigrate(&model.FDSNAccessKey{}); err != nil {
		t.Fatal(err)
	}
	setting := &config.StationFDSNServiceEnabledConfigConstraintImpl{}
	if err := setting.Init(t.Context(), handler); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 18, 12, 0, 0, 0, time.UTC)
	record, key, err := handler.FDSNAccessKeyCreate(context.Background(), "router", nil, now)
	if err != nil {
		t.Fatal(err)
	}
	engine := gin.New()
	if err := Setup(engine, handler, testsupport.NewHardware(), func() time.Time { return now }); err != nil {
		t.Fatal(err)
	}
	var authQueries atomic.Int32
	if err := db.Database.Callback().Query().After("gorm:query").Register("test:fdsn_auth_queries", func(tx *gorm.DB) {
		if tx.Statement.Table == record.GetName(db.GetPrefix()) {
			authQueries.Add(1)
		}
	}); err != nil {
		t.Fatal(err)
	}
	requestMethod := func(method, path, authorization string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, path, nil)
		req.RemoteAddr = "192.0.2.1:4321"
		req.Header.Set("Authorization", authorization)
		req.Header.Set("User-Agent", "fdsn-test")
		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, req)
		return rec
	}
	request := func(path, authorization string) *httptest.ResponseRecorder {
		return requestMethod(http.MethodGet, path, authorization)
	}
	for _, path := range []string{
		"/fdsnws", "/fdsnws/", "/fdsnws/assets/app.js",
		"/fdsnws/station/1/query?network=AB%2AC", "/fdsnws/dataselect/1/query",
	} {
		rec := request(path, "")
		if rec.Code != http.StatusNotImplemented || rec.Header().Get("Content-Type") != fdsnwsgo.MIMEText ||
			rec.Header().Get("WWW-Authenticate") != "" || rec.Header().Get("Cache-Control") != "no-store" {
			t.Fatalf("disabled %s: status = %d", path, rec.Code)
		}
		for _, want := range []string{
			"Error 501: Not Implemented\n\nFDSN service is disabled\n\n",
			"Usage details are not provided by administrator\n\n",
			"Request:\n" + path + "\n\n",
			"Request Submitted:\n" + now.Format(time.RFC3339Nano) + "\n\n",
			"Service version:\n",
		} {
			if !strings.Contains(rec.Body.String(), want) {
				t.Fatalf("disabled %s: missing %q in %s", path, want, rec.Body.String())
			}
		}
	}
	if err := setting.Set(t.Context(), handler, true); err != nil {
		t.Fatal(err)
	}
	for _, page := range []struct {
		path string
		body string
	}{
		{"/fdsnws", fdsnwsRootPage},
		{"/fdsnws/dataselect", dataSelectRootPage},
		{"/fdsnws/dataselect/1", dataSelectVersionRootPage},
		{"/fdsnws/station", stationRootPage},
		{"/fdsnws/station/1", stationVersionRootPage},
	} {
		for _, suffix := range []string{"", "/"} {
			path := page.path + suffix
			for _, authorization := range []string{"", "Bearer invalid"} {
				rec := request(path, authorization)
				if rec.Code != http.StatusOK || rec.Header().Get("WWW-Authenticate") != "" ||
					rec.Header().Get("Content-Type") != "text/html; charset=utf-8" || rec.Body.String() != page.body {
					t.Fatalf("public root %s: status = %d body = %s", path, rec.Code, rec.Body.String())
				}
			}
			if rec := requestMethod(http.MethodPost, path, "Bearer invalid"); rec.Code != http.StatusMethodNotAllowed || rec.Header().Get("WWW-Authenticate") != "" {
				t.Fatalf("POST root %s: status = %d body = %s", path, rec.Code, rec.Body.String())
			}
		}
	}
	for _, path := range []string{
		"/fdsnws/query", "/fdsnws/station/2/query", "/fdsnws/station/1/query/extra",
		"/fdsnws/station/1/query//", "/fdsnws/dataselect/1/queryauth//",
		"/fdsnws/missing.html", "/fdsnws/index.html", "/fdsnws/static/", "/fdsnws/static/index.html",
		"/fdsnws/dataselect/index.html", "/fdsnws/dataselect/1/index.html",
		"/fdsnws/station/index.html", "/fdsnws/station/1/index.html",
		"/fdsnws/assets/", "/fdsnws/assets/missing.js",
		"/fdsnws/assets/../index.html", "/fdsnws/assets-other/app.js",
	} {
		if rec := request(path, "Bearer invalid"); rec.Code != http.StatusNotFound || rec.Header().Get("WWW-Authenticate") != "" {
			t.Fatalf("unknown endpoint %s: status = %d body = %s", path, rec.Code, rec.Body.String())
		}
	}
	if rec := request("/fdsnws/assets", "Bearer invalid"); rec.Code != http.StatusMovedPermanently ||
		rec.Header().Get("Location") != "/fdsnws/assets/" || rec.Header().Get("WWW-Authenticate") != "" {
		t.Fatalf("asset directory redirect: status = %d headers = %v", rec.Code, rec.Header())
	}
	for _, path := range []string{
		"/fdsnws/station/1/version", "/fdsnws/dataselect/1/version",
		"/fdsnws/station/1/application.wadl", "/fdsnws/dataselect/1/application.wadl",
	} {
		for _, authorization := range []string{"", "Bearer invalid"} {
			if rec := request(path, authorization); rec.Code != http.StatusOK || rec.Header().Get("WWW-Authenticate") != "" {
				t.Fatalf("public protocol endpoint %s: status = %d body = %s", path, rec.Code, rec.Body.String())
			}
		}
	}
	for _, asset := range []struct {
		name     string
		mimeType string
	}{
		{"app.js", "javascript"},
		{"app.css", "text/css"},
		{"alpine-csp-3.15.8.min.js", "javascript"},
		{"tailwindcss-4.3.3.js", "javascript"},
	} {
		body, err := webpageAssets.ReadFile("static/assets/" + asset.name)
		if err != nil {
			t.Fatal(err)
		}
		for _, method := range []string{http.MethodGet, http.MethodHead} {
			for _, authorization := range []string{"", "Bearer invalid"} {
				rec := requestMethod(method, "/fdsnws/assets/"+asset.name, authorization)
				if rec.Code != http.StatusOK || rec.Header().Get("WWW-Authenticate") != "" ||
					!strings.Contains(rec.Header().Get("Content-Type"), asset.mimeType) {
					t.Fatalf("public asset %s %s: status = %d headers = %v", method, asset.name, rec.Code, rec.Header())
				}
				if method == http.MethodGet && rec.Body.String() != string(body) {
					t.Fatalf("asset %s differs from the embedded file", asset.name)
				}
				if method == http.MethodHead && rec.Body.Len() != 0 {
					t.Fatalf("HEAD asset %s returned a body", asset.name)
				}
			}
		}
	}
	for _, path := range []string{"/fdsnws/dataselect/1/queryauth", "/fdsnws/dataselect/1/queryauth/?nodata=404"} {
		for _, method := range []string{http.MethodGet, http.MethodPost} {
			for _, authorization := range []string{"", "Bearer invalid", "Bearer " + key} {
				if rec := requestMethod(method, path, authorization); rec.Code != http.StatusNotFound || rec.Header().Get("WWW-Authenticate") != "" {
					t.Fatalf("disabled queryauth %s %s: status = %d body = %s", method, path, rec.Code, rec.Body.String())
				}
			}
		}
	}
	if got := authQueries.Load(); got != 0 {
		t.Fatalf("public requests performed %d authentication queries; want 0", got)
	}
	assertUnauthorized := func(rec *httptest.ResponseRecorder, path, detail string) {
		t.Helper()
		versionPath := strings.SplitN(path, "/query", 2)[0] + "/version"
		version := request(versionPath, "")
		if version.Code != http.StatusOK {
			t.Fatalf("service version %s: status = %d", versionPath, version.Code)
		}
		wantBody := fmt.Sprintf("Error 401: Unauthorized\n\n%s\n\nUsage details are not provided by administrator\n\nRequest:\n%s\n\nRequest Submitted:\n%s\n\nService version:\n%s\n",
			detail, path, now.Format(time.RFC3339Nano), strings.TrimSpace(version.Body.String()))
		if rec.Code != http.StatusUnauthorized || rec.Header().Get("WWW-Authenticate") != bearerChallenge ||
			rec.Header().Get("Content-Type") != fdsnwsgo.MIMEText || rec.Body.String() != wantBody {
			t.Fatalf("unauthorized %s: status = %d headers = %v body = %q; want body = %q", path, rec.Code, rec.Header(), rec.Body.String(), wantBody)
		}
	}
	for _, path := range []string{
		"/fdsnws/station/1/query", "/fdsnws/station/1/query/",
		"/fdsnws/dataselect/1/query", "/fdsnws/dataselect/1/query/?nodata=404",
	} {
		for _, method := range []string{http.MethodGet, http.MethodPost} {
			assertUnauthorized(requestMethod(method, path, ""), path, "FDSN access key is required")
		}
	}
	for _, authorization := range []struct {
		header string
		detail string
	}{
		{"Bearer", "FDSN access key is required"},
		{"Basic " + key, "FDSN access key is required"},
		{"Bearer " + key + " extra", "FDSN access key is required"},
		{"Bearer invalid", "Invalid, expired or revoked FDSN access key"},
		{"Bearer fdsn_" + strings.Repeat("A", 43), "Invalid, expired or revoked FDSN access key"},
	} {
		for _, path := range []string{
			"/fdsnws/station/1/query?network=*&station=AB%2AC",
			"/fdsnws/dataselect/1/query/?nodata=404",
		} {
			assertUnauthorized(request(path, authorization.header), path, authorization.detail)
		}
	}
	if rec := request("/fdsnws/dataselect/1/query", "bearer "+key); rec.Code != http.StatusBadRequest || rec.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("valid status = %d body = %s", rec.Code, rec.Body.String())
	}
	keys, err := handler.FDSNAccessKeyList(context.Background())
	if err != nil || len(keys) != 1 || keys[0].UserIp != "192.0.2.1" || keys[0].UserAgent != "fdsn-test" {
		t.Fatalf("request audit = %+v: %v", keys, err)
	}
	authQueries.Store(0)
	// Missing query parameters must reach protocol validation after one authentication.
	if rec := request("/fdsnws/dataselect/1/query", "Bearer "+key); rec.Code != http.StatusBadRequest {
		t.Fatalf("query status = %d body = %s", rec.Code, rec.Body.String())
	}
	if got := authQueries.Load(); got != 1 {
		t.Fatalf("query performed %d authentication queries; want 1", got)
	}
	if rec := request("/fdsnws/dataselect/1/application.wadl", ""); rec.Code != http.StatusOK || strings.Contains(rec.Body.String(), "queryauth") {
		t.Fatalf("WADL status = %d body = %s", rec.Code, rec.Body.String())
	}
	if err := handler.FDSNAccessKeyRevoke(context.Background(), []string{record.KeyId}, now); err != nil {
		t.Fatal(err)
	}
	expiresAt := now.Add(-time.Minute).UnixMilli()
	_, expiredKey, err := handler.FDSNAccessKeyCreate(context.Background(), "expired", &expiresAt, now.Add(-2*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/fdsnws/station/1/query", "/fdsnws/dataselect/1/query"} {
		for _, invalidKey := range []string{key, expiredKey} {
			assertUnauthorized(request(path, "Bearer "+invalidKey), path, "Invalid, expired or revoked FDSN access key")
		}
	}
	if err := setting.Set(t.Context(), handler, false); err != nil {
		t.Fatal(err)
	}
	if rec := request("/fdsnws/station/1/version", "Bearer "+key); rec.Code != http.StatusNotImplemented {
		t.Fatalf("disabled again status = %d", rec.Code)
	}
}

func TestAccessKeyErrorUsesServiceConfiguration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := station.New(fdsnwsgo.UnimplementedProvider{})
	server, err := fdsnwsgo.NewFDSNWS(fdsnwsgo.UnimplementedProvider{},
		fdsnwsgo.WithModules(fdsnwsgo.ModuleStation),
		fdsnwsgo.WithService(service),
		fdsnwsgo.WithServiceVersion(fdsnwsgo.ModuleStation, "1.1.9"),
		fdsnwsgo.WithDocumentationURL("https://example.org/fdsnws"),
	)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 19, 12, 0, 0, 123, time.FixedZone("test", 8*60*60))
	engine := gin.New()
	engine.GET("/fdsnws/station/1/query", accessKeyMiddleware(nil, func() time.Time { return now }, server, service), func(*gin.Context) {
		t.Error("unauthorized request reached the query handler")
	})
	request := httptest.NewRequest(http.MethodGet, "/fdsnws/station/1/query?network=AB%2AC", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, request)
	if rec.Code != http.StatusUnauthorized || rec.Header().Get("WWW-Authenticate") != bearerChallenge ||
		rec.Header().Get("Content-Type") != fdsnwsgo.MIMEText {
		t.Fatalf("status = %d headers = %v", rec.Code, rec.Header())
	}
	for _, want := range []string{
		"Error 401: Unauthorized\n\nFDSN access key is required\n\n",
		"Usage details are available from https://example.org/fdsnws\n\n",
		"Request:\n/fdsnws/station/1/query?network=AB%2AC\n\n",
		"Request Submitted:\n" + now.UTC().Format(time.RFC3339Nano) + "\n\n",
		"Service version:\n1.1.9\n",
	} {
		if !strings.Contains(rec.Body.String(), want) {
			t.Fatalf("missing %q in %s", want, rec.Body.String())
		}
	}
}

func TestDataSelectServiceRequestContext(t *testing.T) {
	db, handler := testsupport.OpenDAO(t)
	testsupport.InitStation(t, handler)
	start := time.Date(2026, 10, 18, 12, 0, 0, 0, time.UTC)
	testsupport.EnsureSeisTable(t, db, start)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := db.Database.Callback().Row().Before("gorm:row").Register("test:fdsn_cancel_waveform_query", func(tx *gorm.DB) {
		if tx.Statement.Context == ctx {
			// Cancel at the database call so protocol validation cannot mask a missing context binding.
			cancel()
		}
	}); err != nil {
		t.Fatal(err)
	}
	service := newDataSelectService(&provider{handler: handler, device: testsupport.NewHardware()})
	request := fdsnwsgo.Request{
		Method: http.MethodGet,
		Query: url.Values{
			"network": {"*"}, "station": {"*"}, "location": {"*"}, "channel": {"*"},
			"starttime": {start.Format(time.RFC3339)}, "endtime": {start.Add(time.Second).Format(time.RFC3339)},
		},
	}
	if _, err := service.Handle(ctx, "query", request); !errors.Is(err, context.Canceled) {
		t.Fatalf("waveform query ignored cancellation: %v", err)
	}
	if _, err := service.Handle(context.Background(), "query", request); !errors.Is(err, fdsnwsgo.ErrNoData) {
		t.Fatalf("another request inherited canceled context: %v", err)
	}
}
