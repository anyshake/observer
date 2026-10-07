package date_header

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/anyshake/observer/pkg/timesource"
	"github.com/gin-gonic/gin"
)

func TestDateHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	router := gin.New()
	router.Use(New(timesource.New(func() time.Time { return now })))
	router.GET("/", func(ctx *gin.Context) { ctx.Status(http.StatusNoContent) })

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d", recorder.Code)
	}
	if got := recorder.Header().Get("Date"); got != now.Format(http.TimeFormat) {
		t.Fatalf("Date = %q", got)
	}
}
