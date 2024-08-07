package mseed

import (
	"errors"
	"net/http"

	v1 "github.com/anyshake/observer/api/v1"
	"github.com/anyshake/observer/server/response"
	"github.com/anyshake/observer/services/miniseed"
	"github.com/anyshake/observer/utils/logger"
	"github.com/gin-gonic/gin"
)

// @Summary AnyShake Observer MiniSEED data
// @Description List MiniSEED data if action is `show`, or export MiniSEED data in .mseed format if action is `export`
// @Router /mseed [post]
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Produce application/octet-stream
// @Param action formData string true "Action to be performed, either `show` or `export`"
// @Param name formData string false "Name of MiniSEED file to be exported, end with `.mseed`"
// @Failure 400 {object} response.HttpResponse "Failed to list or export MiniSEED data due to invalid request body"
// @Failure 410 {object} response.HttpResponse "Failed to export MiniSEED data due to invalid file name or permission denied"
// @Failure 500 {object} response.HttpResponse "Failed to list or export MiniSEED data due to internal server error"
// @Success 200 {object} response.HttpResponse{data=[]miniSeedFileInfo} "Successfully get list of MiniSEED files"
func (h *MSeed) Register(rg *gin.RouterGroup, resolver *v1.Resolver) error {
	// Get MiniSEED service configuration
	var miniseedService miniseed.MiniSeedService
	serviceConfig, ok := resolver.Config.Services[miniseedService.GetServiceName()]
	if !ok {
		return errors.New("failed to get configuration for MiniSEED service")
	}
	basePath := serviceConfig.(map[string]any)["path"].(string)
	lifeCycle := int(serviceConfig.(map[string]any)["lifecycle"].(float64))

	rg.POST("/mseed", func(c *gin.Context) {
		var binding mseedBinding
		if err := c.ShouldBind(&binding); err != nil {
			logger.GetLogger(h.GetApiName()).Errorln(err)
			response.Error(c, http.StatusBadRequest)
			return
		}

		if binding.Action == "show" {
			fileList, err := h.getMiniSeedList(
				basePath,
				resolver.Config.Stream.Station,
				resolver.Config.Stream.Network,
				lifeCycle,
			)
			if err != nil {
				logger.GetLogger(h.GetApiName()).Errorln(err)
				response.Error(c, http.StatusInternalServerError)
				return
			}

			response.Message(c, "Successfully get MiniSEED file list", fileList)
			return
		}

		if len(binding.Name) == 0 {
			response.Error(c, http.StatusBadRequest)
			return
		}

		fileBytes, err := h.getMiniSeedBytes(basePath, binding.Name)
		if err != nil {
			logger.GetLogger(h.GetApiName()).Errorln(err)
			response.Error(c, http.StatusInternalServerError)
			return
		}

		response.File(c, binding.Name, fileBytes)
	})

	return nil
}
