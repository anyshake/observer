package trace

import (
	"net/http"

	v1 "github.com/anyshake/observer/api/v1"
	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/server/response"
	"github.com/anyshake/observer/utils/logger"
	"github.com/gin-gonic/gin"
)

// @Summary AnyShake Observer event trace
// @Description Get list of earthquake events data source and earthquake events from specified data source
// @Router /trace [post]
// @Accept application/x-www-form-urlencoded
// @Produce application/json
// @Param source formData string true "Use `show` to get available sources first, then choose one and request again to get events"
// @Failure 400 {object} response.HttpResponse "Failed to read earthquake event list due to invalid data source"
// @Failure 500 {object} response.HttpResponse "Failed to read earthquake event list due to failed to read data source"
// @Success 200 {object} response.HttpResponse{data=[]seismicEvent} "Successfully read the list of earthquake events"
func (t *Trace) Register(rg *gin.RouterGroup, resolver *v1.Resolver) error {
	sources := map[string]dataSource{
		"CWA":      &CWA{},
		"HKO":      &HKO{},
		"JMA":      &JMA{},
		"KMA":      &KMA{},
		"CEIC":     &CEIC{},
		"USGS":     &USGS{},
		"INGV":     &INGV{},
		"SCEA_E":   &SCEA_E{},
		"SCEA_B":   &SCEA_B{},
		"CEA_DASE": &CEA_DASE{},
	}

	var explorerDeps *explorer.ExplorerDependency
	err := resolver.Dependency.Invoke(func(deps *explorer.ExplorerDependency) error {
		explorerDeps = deps
		return nil
	})
	if err != nil {
		return err
	}

	rg.POST("/trace", func(c *gin.Context) {
		var binding traceBinding
		if err := c.ShouldBind(&binding); err != nil {
			logger.GetLogger(t.GetApiName()).Errorln(err)
			response.Error(c, http.StatusBadRequest)
			return
		}

		if binding.Source == "show" {
			type availableSources struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			}

			var list []availableSources
			for k, v := range sources {
				name := v.Property()
				list = append(list, availableSources{
					Name:  name,
					Value: k,
				})
			}

			response.Message(c, "Successfully read available data source list", list)
			return
		}

		if err != nil {
			logger.GetLogger(t.GetApiName()).Errorln(err)
			return
		}
		var (
			source, ok = sources[binding.Source]
			latitude   = explorerDeps.Config.GetLatitude()
			longitude  = explorerDeps.Config.GetLongitude()
		)
		if ok {
			events, err := source.List(latitude, longitude)
			if err != nil {
				logger.GetLogger(t.GetApiName()).Errorln(err)
				response.Error(c, http.StatusInternalServerError)
				return
			}

			response.Message(c, "Successfully read the list of earthquake events", sortSeismicEvents(events))
			return
		}

		response.Error(c, http.StatusBadRequest)
	})

	return nil
}
