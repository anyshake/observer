package cors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AllowCros(headers []HttpHeader) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, v := range headers {
			c.Writer.Header().Set(
				v.Header, v.Value,
			)
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
