package mseed

import (
	"net/http"

	"github.com/bclswl0827/observer/app"
	"github.com/bclswl0827/observer/server/response"
	"github.com/gin-gonic/gin"
)

// @Summary Observer MiniSEED data
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
// @Success 200 {object} response.HttpResponse{data=[]MiniSEEDFile} "Successfully get list of MiniSEED files"
func (h *MSeed) RegisterModule(rg *gin.RouterGroup, options *app.ServerOptions) {
	rg.POST("/mseed", func(c *gin.Context) {
		var binding Binding
		if err := c.ShouldBind(&binding); err != nil {
			response.Error(c, http.StatusBadRequest)
			return
		}

		if binding.Action == "show" {
			fileList, err := getMiniSEEDList(options.FeatureOptions.Config)
			if err != nil {
				response.Error(c, http.StatusInternalServerError)
				return
			}

			if len(fileList) == 0 {
				response.Error(c, http.StatusGone)
				return
			}

			response.Message(c, "Successfully get MiniSEED file list", fileList)
			return
		}

		fileBytes, err := getMiniSEEDBytes(options.FeatureOptions.Config, binding.Name)
		if err != nil {
			response.Error(c, http.StatusInternalServerError)
			return
		}

		if len(fileBytes) == 0 {
			response.Error(c, http.StatusGone)
			return
		}

		response.File(c, binding.Name, fileBytes)
	})
}
