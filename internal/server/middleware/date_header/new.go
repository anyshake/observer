package date_header

import (
	"net/http"

	"github.com/anyshake/observer/pkg/timesource"
	"github.com/gin-gonic/gin"
)

func New(timeSource *timesource.Source) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Date", timeSource.Now().UTC().Format(http.TimeFormat))
		c.Next()
	}
}
