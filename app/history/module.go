package history

import (
	"net/http"

	"com.geophone.observer/app"
	"com.geophone.observer/server/response"
	"github.com/gin-gonic/gin"
)

func (s *History) RegisterModule(rg *gin.RouterGroup, options *app.ServerOptions) {
	rg.POST("/history", func(c *gin.Context) {
		var binding Binding
		if err := c.ShouldBind(&binding); err != nil {
			response.ErrorHandler(c, http.StatusBadRequest)
			return
		}

		FilterHistory(c, &binding, options)
	})
}
