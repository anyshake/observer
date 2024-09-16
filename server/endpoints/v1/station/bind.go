package station

import (
	"net/http"

	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/server/response"
	"github.com/anyshake/observer/services"
	"github.com/anyshake/observer/utils/logger"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// @Summary Station Status
// @Description Get Observer station status including system information, memory usage, disk usage, CPU usage, ADC information, geophone information, and location information. This API requires a valid JWT token if the server is in restricted mode.
// @Router /station [get]
// @Produce application/json
// @Security ApiKeyAuth
// @Param Authorization header string false "Bearer JWT token, only required when the server is in restricted mode."
func (s *Station) Bind(rg *gin.RouterGroup, jwtHandler *jwt.GinJWTMiddleware, options *services.Options) error {
	var explorerDeps *explorer.ExplorerDependency
	err := options.Dependency.Invoke(func(deps *explorer.ExplorerDependency) error {
		explorerDeps = deps
		return nil
	})
	if err != nil {
		return err
	}

	var handlerFunc []gin.HandlerFunc
	if options.Config.Server.Restrict {
		handlerFunc = append(handlerFunc, jwtHandler.MiddlewareFunc())
	}
	handlerFunc = append(handlerFunc, func(c *gin.Context) {
		var explorer explorerInfo
		err := explorer.get(options.TimeSource, explorerDeps)
		if err != nil {
			logger.GetLogger(s.GetApiName()).Errorln(err)
			response.Message(c, options.TimeSource, "failed to get explorer information", http.StatusInternalServerError, nil)
			return
		}
		var cpu cpuInfo
		err = cpu.get()
		if err != nil {
			logger.GetLogger(s.GetApiName()).Errorln(err)
			response.Message(c, options.TimeSource, "failed to get CPU information", http.StatusInternalServerError, nil)
			return
		}
		var disk diskInfo
		err = disk.get()
		if err != nil {
			logger.GetLogger(s.GetApiName()).Errorln(err)
			response.Message(c, options.TimeSource, "failed to get disk information", http.StatusInternalServerError, nil)
			return
		}
		var memory memoryInfo
		err = memory.get()
		if err != nil {
			logger.GetLogger(s.GetApiName()).Errorln(err)
			response.Message(c, options.TimeSource, "failed to get memory information", http.StatusInternalServerError, nil)
			return
		}
		var os osInfo
		err = os.get(options.TimeSource)
		if err != nil {
			logger.GetLogger(s.GetApiName()).Errorln(err)
			response.Message(c, options.TimeSource, "failed to get OS information", http.StatusInternalServerError, nil)
			return
		}
		response.Message(c, options.TimeSource, "Successfully read station information", http.StatusOK, stationInfo{
			Station:  options.Config.Station,
			Stream:   options.Config.Stream,
			Sensor:   options.Config.Sensor,
			Explorer: explorer,
			CPU:      cpu,
			Disk:     disk,
			Memory:   memory,
			OS:       os,
		})
	})

	rg.GET("/station", handlerFunc...)
	return nil
}
