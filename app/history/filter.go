package history

import (
	"encoding/json"
	"net/http"
	"strconv"

	"com.geophone.observer/app"
	"com.geophone.observer/features/collector"
	"com.geophone.observer/features/geophone"
	"com.geophone.observer/server/response"
	"github.com/gin-gonic/gin"
)

func FilterHistory(c *gin.Context, b *Binding, options *app.ServerOptions) {
	ctx := options.ConnRedis.Context()
	k, err := options.ConnRedis.Keys(ctx, "*").Result()
	if err != nil {
		response.ErrorHandler(c, http.StatusInternalServerError)
		return
	}

	var acceleration []geophone.Acceleration
	for _, v := range k {
		t, _ := strconv.ParseInt(v, 10, 64)
		if t >= b.Start && t <= b.End {
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
