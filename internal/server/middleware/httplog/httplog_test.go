package httplog

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anyshake/observer/pkg/logger"
	"github.com/gin-gonic/gin"
)

func TestHTTPLog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Init()
	log := logger.GetLogger("httplog")
	router := gin.New()
	router.Use(New(log, "/skip"))
	router.GET("/skip", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })
	router.GET("/ok", func(ctx *gin.Context) { ctx.Status(http.StatusCreated) })
	router.GET("/warn", func(ctx *gin.Context) { ctx.Status(http.StatusBadRequest) })
	router.GET("/error", func(ctx *gin.Context) { ctx.Status(http.StatusInternalServerError) })
	router.GET("/private", func(ctx *gin.Context) {
		_ = ctx.Error(errPrivate("private failure"))
		ctx.Status(http.StatusOK)
	})

	for _, path := range []string{"/skip", "/ok", "/warn", "/error", "/private"} {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
	}
}

type errPrivate string

func (e errPrivate) Error() string { return string(e) }
