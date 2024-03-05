package inventory

import (
	"net/http"

	"github.com/anyshake/observer/app"
	"github.com/anyshake/observer/server/response"
	"github.com/gin-gonic/gin"
)

// @Summary AnyShake Observer station inventory
// @Description Get SeisComP XML inventory, which contains meta data of the station
// @Router /inventory [get]
// @Param format query string false "Format of the inventory, either `json` or `xml`", default is `xml`
// @Produce application/json
// @Success 200 {object} response.HttpResponse{data=string} "Successfully get SeisComP XML inventory"
// @Produce application/xml
func (i *Inventory) RegisterModule(rg *gin.RouterGroup, options *app.ServerOptions) {
	rg.GET("/inventory", func(c *gin.Context) {
		var binding Binding
		if err := c.ShouldBind(&binding); err != nil {
			response.Error(c, http.StatusBadRequest)
			return
		}

		inventory := getInventoryString(options.FeatureOptions.Config, options.FeatureOptions.Status)
		if binding.Format == "json" {
			response.Message(c, "Successfully get SeisComP XML inventory", inventory)
			return
		}

		c.Data(http.StatusOK, "application/xml", []byte(inventory))
	})
}
