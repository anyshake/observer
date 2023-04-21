package history

import (
	"encoding/json"
	"net/http"

	"com.geophone.observer/app"
	"com.geophone.observer/common/postgres"
	"com.geophone.observer/features/geophone"
	"com.geophone.observer/server/response"
	"github.com/gin-gonic/gin"
)

func FilterHistory(c *gin.Context, b *Binding, options *app.ServerOptions) {
	data, err := postgres.SelectData(options.ConnPostgres, b.Timestamp)
	if err != nil {
		response.ErrorHandler(c, http.StatusInternalServerError)
		return
	}

	if len(data) == 0 {
		response.ErrorHandler(c, http.StatusBadRequest)
		return
	}

	var acceleration []geophone.Acceleration
	for _, v := range data {
		var acc []geophone.Acceleration
		err := json.Unmarshal([]byte(v["data"].(string)), &acc)
		if err != nil {
			response.ErrorHandler(c, http.StatusInternalServerError)
			return
		}

		acceleration = append(acceleration, acc...)
	}

	c.JSON(http.StatusOK, response.MessageHandler(c, "筛选出如下加速度数据", acceleration))
}
