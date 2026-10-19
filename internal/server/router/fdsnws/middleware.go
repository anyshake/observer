package fdsnws

import (
	"net/http"
	"strings"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
	fdsnwsgo "github.com/bclswl0827/fdsnws-go"
	"github.com/gin-gonic/gin"
)

func serviceMiddleware(handler *action.Handler, server *fdsnwsgo.FDSNWS) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "no-store")
		enabled, err := (&config.StationFDSNServiceEnabledConfigConstraintImpl{}).Get(ctx.Request.Context(), handler)
		status := http.StatusNotImplemented
		detail := "FDSN service is disabled"
		if err != nil {
			_ = ctx.Error(err)
			status = http.StatusServiceUnavailable
			detail = "FDSN service configuration is unavailable"
		} else if enabled.(bool) {
			// Bound bulk POST requests before the protocol adapter reads their body.
			ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 1<<20)
			ctx.Next()
			return
		}

		module, _, _ := strings.Cut(strings.TrimPrefix(ctx.Request.URL.Path, basePath+"/"), "/")
		response := server.ErrorResponse(fdsnwsgo.Request{
			URL: ctx.Request.URL.String(),
		}, fdsnwsgo.Module(module), status, http.StatusText(status), detail)
		ctx.Abort()
		if err := fdsnwsgo.WriteHTTPResponse(ctx.Writer, response); err != nil {
			_ = ctx.Error(err)
		}
	}
}
