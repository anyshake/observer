package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestResponseWriters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/data", func(ctx *gin.Context) {
		Data(ctx, http.StatusCreated, "ok", map[string]int{"n": 1})
	})
	router.GET("/error", func(ctx *gin.Context) {
		Error(ctx, http.StatusBadRequest, "bad")
	})
	router.GET("/blob", func(ctx *gin.Context) {
		Blob(ctx, "trace.txt", "text/plain", []byte("wave"))
	})

	for _, tc := range []struct {
		path       string
		status     int
		wantHeader string
		wantBody   string
	}{
		{"/data", http.StatusOK, "", `"n":1`},
		{"/error", http.StatusBadRequest, "", `"error":true`},
		{"/blob", http.StatusOK, "attachment; filename=trace.txt", "wave"},
	} {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if recorder.Code != tc.status {
			t.Fatalf("%s status = %d", tc.path, recorder.Code)
		}
		if tc.wantHeader != "" && recorder.Header().Get("Content-Disposition") != tc.wantHeader {
			t.Fatalf("%s disposition = %q", tc.path, recorder.Header().Get("Content-Disposition"))
		}
		if body := recorder.Body.String(); !strings.Contains(body, tc.wantBody) {
			t.Fatalf("%s body = %s", tc.path, body)
		}
	}

	var parsed Response
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/data", nil))
	if err := json.Unmarshal(recorder.Body.Bytes(), &parsed); err != nil || parsed.Error || parsed.Message != "ok" {
		t.Fatalf("parsed data response = %+v, %v", parsed, err)
	}
}
