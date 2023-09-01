package trace

import (
	"net/http"
	"time"

	"github.com/bclswl0827/observer/app"
	"github.com/bclswl0827/observer/server/middleware/limit"
	"github.com/bclswl0827/observer/server/response"
	"github.com/gin-gonic/gin"
)

func (t *Trace) RegisterModule(rg *gin.RouterGroup, options *app.ServerOptions) {
	sources := []DataSource{
		&USGS{}, &JMA{}, &CWB{}, &HKO{},
		&CEIC{}, &SCEA_E{}, &SCEA_B{},
	}

	rg.Use(limit.RateLimit(time.Second, CAPACITY, CAPACITY))
	rg.POST("/trace", func(c *gin.Context) {
		var binding Binding
		if err := c.ShouldBind(&binding); err != nil {
			response.Error(c, http.StatusBadRequest)
			return
		}

		if binding.Source == "show" {
			type availableSources struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			}

			var list []availableSources
			for _, v := range sources {
				name, value := v.Property()
				list = append(list, availableSources{
					Name:  name,
					Value: value,
				})
			}

			response.Message(c, "成功取得可用数据源列表", list)
			return
		}

		var (
			latitude  = options.FeatureOptions.Config.Station.Latitude
			longitude = options.FeatureOptions.Config.Station.Longitude
		)
		for _, v := range sources {
			_, value := v.Property()
			if value == binding.Source {
				events, err := v.List(latitude, longitude)
				if err != nil {
					response.Error(c, http.StatusInternalServerError)
					return
				}

				sortByTimestamp(events)
				response.Message(c, "成功取得地震列表数据", events)
				return
			}
		}

		response.Error(c, http.StatusBadRequest)
	})
}
