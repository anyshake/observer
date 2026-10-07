package recovery

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anyshake/observer/pkg/logger"
	"github.com/gin-gonic/gin"
)

func TestRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Init()
	router := gin.New()
	router.Use(New(logger.GetLogger("recovery")))
	router.GET("/panic", func(*gin.Context) { panic("boom") })
	router.GET("/ok", func(ctx *gin.Context) { ctx.Status(http.StatusNoContent) })

	panicked := httptest.NewRecorder()
	router.ServeHTTP(panicked, httptest.NewRequest(http.MethodGet, "/panic", nil))
	if panicked.Code != http.StatusInternalServerError {
		t.Fatalf("panic status = %d", panicked.Code)
	}

	ok := httptest.NewRecorder()
	router.ServeHTTP(ok, httptest.NewRequest(http.MethodGet, "/ok", nil))
	if ok.Code != http.StatusNoContent {
		t.Fatalf("ok status = %d", ok.Code)
	}
}
