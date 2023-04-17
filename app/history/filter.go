package history

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"com.geophone.observer/app"
	"com.geophone.observer/features/collector"
	"com.geophone.observer/features/geophone"
	"com.geophone.observer/server/response"
	"github.com/gin-gonic/gin"
)

func FilterHistory(c *gin.Context, b *Binding, options *app.ServerOptions) {
	ctx := options.ConnRedis.Context()

	var keys []string
	iter := options.ConnRedis.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		response.ErrorHandler(c, http.StatusInternalServerError)
		return
	}

	var acceleration []geophone.Acceleration
	for _, v := range keys {
		t, _ := strconv.ParseInt(v, 10, 64)
		if t >= b.Timestamp && t <= time.UnixMilli(b.Timestamp).Add(time.Minute+2*time.Second).UnixMilli() {
			value, err := options.ConnRedis.Get(ctx, v).Result()
			if err != nil {
				response.ErrorHandler(c, http.StatusInternalServerError)
				return
			}

			var message collector.Message
			json.Unmarshal([]byte(value), &message)
			for _, _v := range message.Acceleration {
				acceleration = append(acceleration, _v)
			}
		}
	}

	if len(acceleration) == 0 {
		response.ErrorHandler(c, http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, response.MessageHandler(c, "筛选出如下加速度数据", acceleration))
}
