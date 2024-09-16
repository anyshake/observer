package history

import (
	"errors"
	"net/http"

	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/server/response"
	"github.com/anyshake/observer/services"
	"github.com/anyshake/observer/utils/logger"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary Waveform History
// @Description Get seismic waveform data from database in specified time range, channel and format. This API supports 1 hour of maximum duration of the waveform data to be queried. This API requires a valid JWT token if the server is in restricted mode.
// @Router /history [post]
// @Produce application/json
// @Produce application/octet-stream
// @Security ApiKeyAuth
// @Param start_time formData int true "Start time of the waveform to be queried, unix timestamp format in milliseconds."
// @Param end_time formData int true "End time of the waveform to be queried, unix timestamp format in milliseconds."
// @Param format formData string true "Set output format of the waveform data, available options are `json`, `sac`, and `miniseed`."
// @Param channel formData string false "Channel of the waveform, available options are `Z`, `E` or `N` (in uppercase), only reuqired when output format is set to `sac` and `miniseed`."
// @Param Authorization header string false "Bearer JWT token, only required when the server is in restricted mode."
func (h *History) Bind(rg *gin.RouterGroup, jwtHandler *jwt.GinJWTMiddleware, options *services.Options) error {
	var handlerFunc []gin.HandlerFunc
	if options.Config.Server.Restrict {
		handlerFunc = append(handlerFunc, jwtHandler.MiddlewareFunc())
	}
	handlerFunc = append(handlerFunc, func(c *gin.Context) {
		var req request
		err := c.ShouldBind(&req)
		if err != nil {
			logger.GetLogger(h.GetApiName()).Errorln(err)
			response.Message(c, options.TimeSource, "request body is not valid", http.StatusBadRequest, nil)
			return
		}

		// Filter out the history data
		resultArr, err := h.filter(req.StartTime, req.EndTime, JSON_MAX_DURATION, options.Database)
		if err != nil {
			logger.GetLogger(h.GetApiName()).Errorln(err)
			response.Message(c, options.TimeSource, "no data available for the given time range", http.StatusGone, nil)
			return
		}

		switch req.Format {
		case "json":
			response.Message(c, options.TimeSource, "successfully exported the waveform data in JSON format", http.StatusOK, resultArr)
			return
		case "sac":
			if req.Channel != explorer.EXPLORER_CHANNEL_CODE_Z &&
				req.Channel != explorer.EXPLORER_CHANNEL_CODE_E &&
				req.Channel != explorer.EXPLORER_CHANNEL_CODE_N {
				logger.GetLogger(h.GetApiName()).Errorln(errors.New("no channel was selected in SAC format"))
				response.Message(c, options.TimeSource, "failed to export waveform data due to invalid channel", http.StatusBadRequest, nil)
				return
			}

			fileName, dataBytes, err := h.handleSAC(
				resultArr,
				options.Config.Stream.Station,
				options.Config.Stream.Network,
				options.Config.Stream.Location,
				options.Config.Stream.Channel,
				req.Channel,
			)
			if err != nil {
				logger.GetLogger(h.GetApiName()).Errorln(err)
				response.Message(c, options.TimeSource, err.Error(), http.StatusInternalServerError, nil)
				return
			}

			response.File(c, fileName, dataBytes)
			return
		case "miniseed":
			if req.Channel != explorer.EXPLORER_CHANNEL_CODE_Z &&
				req.Channel != explorer.EXPLORER_CHANNEL_CODE_E &&
				req.Channel != explorer.EXPLORER_CHANNEL_CODE_N {
				logger.GetLogger(h.GetApiName()).Errorln(errors.New("no channel was selected in MiniSEED format"))
				response.Message(c, options.TimeSource, "failed to export waveform data due to invalid channel", http.StatusBadRequest, nil)
				return
			}

			fileName, dataBytes, err := h.handleMiniSEED(
				resultArr,
				options.Config.Stream.Station,
				options.Config.Stream.Network,
				options.Config.Stream.Location,
				options.Config.Stream.Channel,
				req.Channel,
			)
			if err != nil {
				logger.GetLogger(h.GetApiName()).Errorln(err)
				response.Message(c, options.TimeSource, err.Error(), http.StatusInternalServerError, nil)
				return
			}

			response.File(c, fileName, dataBytes)
			return
		}
	})

	rg.POST("/history", handlerFunc...)
	return nil
}
