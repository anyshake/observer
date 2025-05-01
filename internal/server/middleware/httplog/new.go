package httplog

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func New(logger logrus.FieldLogger, notLogged ...string) gin.HandlerFunc {
	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, p := range notLogged {
			skip[p] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()

		if _, ok := skip[path]; ok {
			return
		}

		if len(c.Errors) > 0 {
			logger.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("%s - \"%s %s\" %d \"%s\" (%d ms)", clientIP, c.Request.Method, path, statusCode, clientUserAgent, latency)
			if statusCode >= http.StatusInternalServerError {
				logger.Error(msg)
			} else if statusCode >= http.StatusBadRequest {
				logger.Warn(msg)
			} else {
				logger.Info(msg)
			}
		}
	}
}
