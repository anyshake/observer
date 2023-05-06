package trace

import (
	"net/http"

	"com.geophone.observer/app"
	"com.geophone.observer/server/response"
	"github.com/gin-gonic/gin"
)

func (t *Trace) RegisterModule(rg *gin.RouterGroup, options *app.ServerOptions) {
	sources := []DataSource{
		&HKO{}, &CENC{}, &CWB{}, &USGS{},
	}

	rg.POST("/trace", func(c *gin.Context) {
		var binding Binding
		if err := c.ShouldBind(&binding); err != nil {
			response.ErrorHandler(c, http.StatusBadRequest)
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

			c.JSON(http.StatusOK, response.MessageHandler(c, "成功取得可用数据源列表", list))
			return
		}

		for _, v := range sources {
			_, value := v.Property()
			if value == binding.Source {
				events, err := v.List(options.Message.Latitude, options.Message.Longitude)
				if err != nil {
					response.ErrorHandler(c, http.StatusInternalServerError)
					return
				}

				c.JSON(http.StatusOK, response.MessageHandler(c, "成功取得地震列表数据", events))
				return
			}
		}

		response.ErrorHandler(c, http.StatusBadRequest)
	})
}
