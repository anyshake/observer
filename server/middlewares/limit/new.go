package limit

import (
	"net/http"
	"time"

	"github.com/anyshake/observer/server/response"
	"github.com/anyshake/observer/utils/timesource"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

func New(timeSource *timesource.Source, interval time.Duration, capacity, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(interval, capacity, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			response.Message(c, timeSource, "rate limit exceeded", http.StatusTooManyRequests, nil)
			return
		}

		c.Next()
	}
}
