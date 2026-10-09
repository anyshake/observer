package export

import (
	"fmt"
	"net/http"
	"time"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/server/response"
	"github.com/bclswl0827/mseedio"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

func Setup(routerGroup *gin.RouterGroup, actionHandler *action.Handler, hardware hardware.IHardware, jwtMiddleware gin.HandlerFunc) {
	var (
		stationCodeConfig  = config.StationStationCodeConfigConstraintImpl{}
		locationCodeConfig = config.StationLocationCodeConfigConstraintImpl{}
		networkCodeConfig  = config.StationNetworkCodeConfigConstraintImpl{}
	)
	formats := map[string]seismicDataEncoder{
		"sac": &seismicDataEncoderSacImpl{
			actionHandler:      actionHandler,
			stationCodeConfig:  stationCodeConfig,
			locationCodeConfig: locationCodeConfig,
			networkCodeConfig:  networkCodeConfig,
		},
		"mseed_int32": &seismicDataEncoderMseedImpl{
			encodeType:         mseedio.INT32,
			name:               "MiniSEED (INT32)",
			actionHandler:      actionHandler,
			stationCodeConfig:  stationCodeConfig,
			locationCodeConfig: locationCodeConfig,
			networkCodeConfig:  networkCodeConfig,
		},
		"mseed_steim1": &seismicDataEncoderMseedImpl{
			encodeType:         mseedio.STEIM1,
			name:               "MiniSEED (STEIM1)",
			actionHandler:      actionHandler,
			stationCodeConfig:  stationCodeConfig,
			locationCodeConfig: locationCodeConfig,
			networkCodeConfig:  networkCodeConfig,
		},
		"mseed_steim2": &seismicDataEncoderMseedImpl{
			encodeType:         mseedio.STEIM2,
			name:               "MiniSEED (STEIM2)",
			actionHandler:      actionHandler,
			stationCodeConfig:  stationCodeConfig,
			locationCodeConfig: locationCodeConfig,
			networkCodeConfig:  networkCodeConfig,
		},
		"txt": &seismicDataEncoderTxtImpl{
			actionHandler:      actionHandler,
			stationCodeConfig:  stationCodeConfig,
			locationCodeConfig: locationCodeConfig,
			networkCodeConfig:  networkCodeConfig,
		},
		"wav": &seismicDataEncoderWavImpl{
			outputSampleRate:   44100,
			actionHandler:      actionHandler,
			stationCodeConfig:  stationCodeConfig,
			locationCodeConfig: locationCodeConfig,
			networkCodeConfig:  networkCodeConfig,
		},
	}

	routerGroup.GET("/export", jwtMiddleware, func(ctx *gin.Context) {
		dataFormatMap := make(map[string]string)
		for k, v := range formats {
			dataFormatMap[k] = v.GetName()
		}
		hardwareConfig := hardware.GetConfig()
		channelCodes := append([]string{"*"}, hardwareConfig.GetChannelCodes()...)
		response.Data(ctx, http.StatusOK, "data formats", gin.H{
			"data_format":  dataFormatMap,
			"channel_code": channelCodes,
		})
	})
	routerGroup.POST("/export", jwtMiddleware, func(ctx *gin.Context) {
		var requestModel struct {
			StartTime   int64  `form:"start_time" json:"start_time" xml:"start_time" binding:"required"`
			EndTime     int64  `form:"end_time" json:"end_time" xml:"end_time" binding:"required"`
			ChannelCode string `form:"channel_code" json:"channel_code" xml:"channel_code" binding:"required"`
			DataFormat  string `form:"data_format" json:"data_format" xml:"data_format" binding:"required"`
		}
		if err := ctx.ShouldBind(&requestModel); err != nil {
			response.Error(ctx, http.StatusBadRequest, "request body is not valid")
			return
		}

		hardwareConfig := hardware.GetConfig()
		if requestModel.ChannelCode != "*" && !lo.Contains(hardwareConfig.GetChannelCodes(), requestModel.ChannelCode) {
			response.Error(ctx, http.StatusBadRequest, fmt.Sprintf("channel code %s was not found in hardware config", requestModel.ChannelCode))
			return
		}

		encoder, ok := formats[requestModel.DataFormat]
		if !ok {
			response.Error(ctx, http.StatusBadRequest, fmt.Sprintf("unknown data format type: %s", requestModel.DataFormat))
			return
		}

		if requestModel.ChannelCode == "*" {
			switch requestModel.DataFormat {
			case "sac", "txt", "wav":
				response.Error(ctx, http.StatusBadRequest, fmt.Sprintf("data format %s does not support exporting all channels", requestModel.DataFormat))
				return
			}
		}

		startTime, endTime := requestModel.StartTime, requestModel.EndTime
		startTimestamp := time.UnixMilli(startTime)
		endTimestamp := time.UnixMilli(endTime)
		recordFound := false
		records := seismicRecordIterator(func(callback func(model.SeisRecord) error) error {
			return actionHandler.SeisRecordsQueryEachContext(ctx.Request.Context(), startTimestamp, endTimestamp, func(record model.SeisRecord) error {
				recordFound = true
				return callback(record)
			})
		})
		dataBytes, err := encoder.Encode(records, requestModel.ChannelCode)
		if err != nil {
			err = fmt.Errorf("failed to encode seismic records: %w", err)
			response.Error(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		if !recordFound {
			response.Error(ctx, http.StatusNotFound, "no seis records found in given time range")
			return
		}

		if len(dataBytes) > 0 {
			fileName, err := encoder.GetFileName(startTimestamp, requestModel.ChannelCode)
			if err != nil {
				err = fmt.Errorf("failed to get file name: %w", err)
				response.Error(ctx, http.StatusInternalServerError, err.Error())
				return
			}
			response.Blob(ctx, fileName, "application/octet-stream", dataBytes)
			return
		}

		response.Error(ctx, http.StatusBadRequest, fmt.Sprintf("unknown data type: %s", requestModel.DataFormat))
	})
}
