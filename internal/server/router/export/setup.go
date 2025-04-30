package export

import (
	"fmt"
	"net/http"
	"time"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/server/response"
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
		"mseed": &seismicDataEncoderMseedImpl{
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
		response.Data(ctx, http.StatusOK, "data formats", gin.H{
			"data_format":  dataFormatMap,
			"channel_code": hardwareConfig.GetChannelCodes(),
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
		if !lo.Contains(hardwareConfig.GetChannelCodes(), requestModel.ChannelCode) {
			response.Error(ctx, http.StatusBadRequest, fmt.Sprintf("channel code %s was not found in hardware config", requestModel.ChannelCode))
			return
		}

		encoder, ok := formats[requestModel.DataFormat]
		if !ok {
			response.Error(ctx, http.StatusBadRequest, fmt.Sprintf("unknown data format type: %s", requestModel.DataFormat))
			return
		}

		startTime, endTime := requestModel.StartTime, requestModel.EndTime
		seisRecords, err := actionHandler.SeisRecordsQuery(time.UnixMilli(startTime), time.UnixMilli(endTime))
		if err != nil {
			err = fmt.Errorf("failed to query seis records: %w", err)
			response.Error(ctx, http.StatusInternalServerError, err.Error())
			return
		}

		if len(seisRecords) == 0 {
			response.Error(ctx, http.StatusNotFound, "no seis records found in given time range")
			return
		}

		dataBytes, err := encoder.Encode(seisRecords, requestModel.ChannelCode)
		if err != nil {
			err = fmt.Errorf("failed to encode seismic records: %w", err)
			response.Error(ctx, http.StatusInternalServerError, err.Error())
			return
		}

		if len(dataBytes) > 0 {
			fileName, err := encoder.GetFileName(time.UnixMilli(startTime), requestModel.ChannelCode)
			if err != nil {
				err = fmt.Errorf("failed to get file name: %w", err)
				response.Error(ctx, http.StatusInternalServerError, err.Error())
				return
			}
			response.Blob(ctx, fileName, dataBytes)
			return
		}

		response.Error(ctx, http.StatusBadRequest, fmt.Sprintf("unknown data type: %s", requestModel.DataFormat))
	})
}
