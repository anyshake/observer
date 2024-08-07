package history

import (
	"fmt"
	"net/http"

	v1 "github.com/anyshake/observer/api/v1"
	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/server/response"
	"github.com/anyshake/observer/utils/logger"
	"github.com/gin-gonic/gin"
)

// @Summary AnyShake Observer waveform history
// @Description Get waveform count data in specified time range, channel and format, the maximum duration of the waveform data to be exported is 24 hours for JSON and 1 hour for SAC
// @Router /history [post]
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Produce application/octet-stream
// @Param start_time formData int true "Start timestamp of the waveform data to be queried, in milliseconds (unix timestamp)"
// @Param end_time formData int true "End timestamp of the waveform data to be queried, in milliseconds (unix timestamp)"
// @Param format formData string true "Format of the waveform data to be queried, `json` or `sac`"
// @Param channel formData string false "Channel of the waveform, `Z`, `E` or `N`, reuqired when format is `sac`"
// @Failure 400 {object} response.HttpResponse "Failed to export waveform data due to invalid format or channel"
// @Failure 410 {object} response.HttpResponse "Failed to export waveform data due to no data available"
// @Failure 500 {object} response.HttpResponse "Failed to export waveform data due to failed to read data source"
// @Success 200 {object} response.HttpResponse{data=[]explorer.ExplorerData} "Successfully exported the waveform data"
func (h *History) Register(rg *gin.RouterGroup, resolver *v1.Resolver) error {
	rg.POST("/history", func(c *gin.Context) {
		var binding historyBinding
		if err := c.ShouldBind(&binding); err != nil {
			logger.GetLogger(h.GetApiName()).Errorln(err)
			response.Error(c, http.StatusBadRequest)
			return
		}

		switch binding.Format {
		case "json":
			result, err := h.filterHistory(binding.StartTime, binding.EndTime, JSON_MAX_DURATION, resolver)
			if err != nil {
				logger.GetLogger(h.GetApiName()).Errorln(err)
				response.Error(c, http.StatusGone)
				return
			}
			response.Message(c, "The waveform data was successfully filtered", result)
			return
		case "sac":
			result, err := h.filterHistory(binding.StartTime, binding.EndTime, SAC_MAX_DURATION, resolver)
			if err != nil {
				logger.GetLogger(h.GetApiName()).Errorln(err)
				response.Error(c, http.StatusGone)
				return
			}
			if binding.Channel != explorer.EXPLORER_CHANNEL_CODE_Z &&
				binding.Channel != explorer.EXPLORER_CHANNEL_CODE_E &&
				binding.Channel != explorer.EXPLORER_CHANNEL_CODE_N {
				err := fmt.Errorf("no channel was selected")
				logger.GetLogger(h.GetApiName()).Errorln(err)
				response.Error(c, http.StatusBadRequest)
				return
			}
			fileName, dataBytes, err := h.getSACBytes(
				result,
				resolver.Config.Stream.Station,
				resolver.Config.Stream.Network,
				resolver.Config.Stream.Location,
				resolver.Config.Stream.Channel,
				binding.Channel,
			)
			if err != nil {
				logger.GetLogger(h.GetApiName()).Errorln(err)
				response.Error(c, http.StatusInternalServerError)
				return
			}

			response.File(c, fileName, dataBytes)
			return
		}

		response.Error(c, http.StatusBadRequest)
	})

	return nil
}
