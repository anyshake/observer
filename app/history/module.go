package history

import (
	"net/http"

	"com.geophone.observer/app"
	"com.geophone.observer/server/response"
	"github.com/gin-gonic/gin"
)

func (h *History) RegisterModule(rg *gin.RouterGroup, options *app.ServerOptions) {
	rg.POST("/history", func(c *gin.Context) {
		var binding Binding
		if err := c.ShouldBind(&binding); err != nil {
			response.ErrorHandler(c, http.StatusBadRequest)
			return
		}

		data, err := FilterHistory(binding.Timestamp, options)
		if err != nil {
			response.ErrorHandler(c, http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, response.MessageHandler(c, "筛选出如下加速度数据", data))
	})
}
