package station

import (
	"net/http"

	v1 "github.com/anyshake/observer/api/v1"
	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/server/response"
	"github.com/anyshake/observer/utils/logger"
	"github.com/gin-gonic/gin"
)

// @Summary AnyShake Observer station status
// @Description Get Observer station status including system information, memory usage, disk usage, CPU usage, ADC information, geophone information, and location information
// @Router /station [get]
// @Produce application/json
// @Success 200 {object} response.HttpResponse{data=stationInfo} "Successfully read station information"
func (s *Station) Register(rg *gin.RouterGroup, resolver *v1.Resolver) error {
	var explorerDeps *explorer.ExplorerDependency
	err := resolver.Dependency.Invoke(func(deps *explorer.ExplorerDependency) error {
		explorerDeps = deps
		return nil
	})
	if err != nil {
		return err
	}

	rg.GET("/station", func(c *gin.Context) {
		var explorer explorerInfo
		err := explorer.get(resolver.TimeSource, explorerDeps)
		if err != nil {
			logger.GetLogger(s.GetApiName()).Errorln(err)
			response.Error(c, http.StatusInternalServerError)
			return
		}
		var cpu cpuInfo
		err = cpu.get()
		if err != nil {
			logger.GetLogger(s.GetApiName()).Errorln(err)
			response.Error(c, http.StatusInternalServerError)
			return
		}
		var disk diskInfo
		err = disk.get()
		if err != nil {
			logger.GetLogger(s.GetApiName()).Errorln(err)
			response.Error(c, http.StatusInternalServerError)
			return
		}
		var memory memoryInfo
		err = memory.get()
		if err != nil {
			logger.GetLogger(s.GetApiName()).Errorln(err)
			response.Error(c, http.StatusInternalServerError)
			return
		}
		var os osInfo
		err = os.get(resolver.TimeSource)
		if err != nil {
			logger.GetLogger(s.GetApiName()).Errorln(err)
			response.Error(c, http.StatusInternalServerError)
			return
		}
		response.Message(c, "Successfully read station information", stationInfo{
			Station:  resolver.Config.Station,
			Stream:   resolver.Config.Stream,
			Sensor:   resolver.Config.Sensor,
			Explorer: explorer,
			CPU:      cpu,
			Disk:     disk,
			Memory:   memory,
			OS:       os,
		})
	})

	return nil
}
