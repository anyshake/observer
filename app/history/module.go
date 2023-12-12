package history

import (
	"net/http"

	"github.com/anyshake/observer/app"
	"github.com/anyshake/observer/server/response"
	"github.com/gin-gonic/gin"
)

// @Summary Observer waveform history
// @Description Get waveform count data in specified time range, channel and format
// @Router /history [post]
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Produce application/octet-stream
// @Param start formData int true "Start timestamp of the waveform data to be queried, in milliseconds"
// @Param end formData int true "End timestamp of the waveform data to be queried, in milliseconds"
// @Param format formData string true "Format of the waveform data to be queried, `json` or `sac`"
// @Param channel formData string false "Channel of the waveform, `EHZ`, `EHE` or `EHN`, reuqired when format is `sac`"
// @Failure 400 {object} response.HttpResponse "Failed to export waveform data due to invalid format or channel"
// @Failure 410 {object} response.HttpResponse "Failed to export waveform data due to no data available"
// @Failure 500 {object} response.HttpResponse "Failed to export waveform data due to failed to read data source"
// @Success 200 {object} response.HttpResponse{data=[]publisher.Geophone} "Successfully exported the waveform data"
func (h *History) RegisterModule(rg *gin.RouterGroup, options *app.ServerOptions) {
	rg.POST("/history", func(c *gin.Context) {
		var binding Binding
		if err := c.ShouldBind(&binding); err != nil {
			response.Error(c, http.StatusBadRequest)
			return
		}

		data, err := filterHistory(binding.Start, binding.End, options)
		if err != nil {
			response.Error(c, http.StatusGone)
			return
		}

		switch binding.Format {
		case "json":
			response.Message(c, "The specified waveform data was successfully filtered", data)
			return
		case "sac":
			fileName, dataBytes, err := getSACBytes(data, binding.Channel, options)
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
