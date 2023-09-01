package limit

import (
	"net/http"
	"time"

	"github.com/bclswl0827/observer/server/response"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

func RateLimit(interval time.Duration, capacity, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(interval, capacity, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			response.Error(c, http.StatusTooManyRequests)
			return
		}

		c.Next()
	}
}
