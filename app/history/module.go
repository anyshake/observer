package history

import (
	"net/http"
	"time"

	"github.com/bclswl0827/observer/app"
	"github.com/bclswl0827/observer/server/middleware/limit"
	"github.com/bclswl0827/observer/server/response"
	"github.com/gin-gonic/gin"
)

func (h *History) RegisterModule(rg *gin.RouterGroup, options *app.ServerOptions) {
	rg.Use(limit.RateLimit(time.Second, CAPACITY, CAPACITY))
	rg.POST("/history", func(c *gin.Context) {
		var binding Binding
		if err := c.ShouldBind(&binding); err != nil {
			response.Error(c, http.StatusBadRequest)
			return
		}

		data, err := FilterHistory(binding.Start, binding.End, options)
		if err != nil {
			response.Error(c, http.StatusGone)
			return
		}

		switch binding.Format {
		case "json":
			response.Message(c, "The specified waveform data was successfully filtered", data)
			return
		case "sac":
			fileName, dataBytes, err := GetSACBytes(data, binding.Channel, options)
			if err != nil {
				response.Error(c, http.StatusInternalServerError)
				return
			}

			response.File(c, fileName, dataBytes)
			return
		}

		response.Error(c, http.StatusBadRequest)
	})
}
