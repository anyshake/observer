package history

import (
	"net/http"
	"time"

	"github.com/bclswl0827/observer/app"
	"github.com/bclswl0827/observer/server/middleware/limit"
	"github.com/bclswl0827/observer/server/response"
	"github.com/gin-gonic/gin"
)

// @Summary Observer waveform history
// @Description Get earthquake events data source list and earthquake event list from data source
// @Router /history [post]
// @Accept application/json
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Produce application/octet-stream
// @Param start body int true "Start timestamp of the waveform data to be queried, in milliseconds"
// @Param end body int true "End timestamp of the waveform data to be queried, in milliseconds"
// @Param format body string true "Format of the waveform data to be exported, `json` or `sac`"
// @Param channel body string true "Channel of the waveform data to be queried, `EHZ`, `EHE` or `EHN`"
// @Failure 400 {object} response.HttpResponse "Failed to export waveform data due to invalid format or channel"
// @Failure 410 {object} response.HttpResponse "Failed to export waveform data due to no data available"
// @Failure 500 {object} response.HttpResponse "Failed to export waveform data due to failed to read data source"
// @Success 200 {object} response.HttpResponse{data=[]publisher.Geophone} "Successfully exported the waveform data"
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
