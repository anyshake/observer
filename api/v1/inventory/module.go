package inventory

import (
	"net/http"

	v1 "github.com/anyshake/observer/api/v1"
	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/server/response"
	"github.com/anyshake/observer/utils/logger"
	"github.com/gin-gonic/gin"
)

// @Summary AnyShake Observer station inventory
// @Description Get SeisComP XML inventory, which contains meta data of the station
// @Router /inventory [get]
// @Param format query string false "Format of the inventory, either `json` or `xml`", default is `xml`
// @Produce application/json
// @Success 200 {object} response.HttpResponse{data=string} "Successfully get SeisComP XML inventory"
// @Produce application/xml
func (i *Inventory) Register(rg *gin.RouterGroup, resolver *v1.Resolver) error {
	var explorerDeps *explorer.ExplorerDependency
	err := resolver.Dependency.Invoke(func(deps *explorer.ExplorerDependency) error {
		explorerDeps = deps
		return nil
	})
	if err != nil {
		return err
	}

	rg.GET("/inventory", func(c *gin.Context) {
		var binding inventoryBinding
		if err := c.ShouldBind(&binding); err != nil {
			logger.GetLogger(i.GetApiName()).Errorln(err)
			response.Error(c, http.StatusBadRequest)
			return
		}

		inventory := i.getInventoryString(resolver.Config, explorerDeps)
		if binding.Format == "json" {
			response.Message(c, "Successfully get SeisComP XML inventory", inventory)
			return
		}

		c.Data(http.StatusOK, "application/xml", []byte(inventory))
	})

	return nil
}
