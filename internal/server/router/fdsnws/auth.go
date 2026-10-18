package fdsnws

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	fdsnwsgo "github.com/bclswl0827/fdsnws-go"
	"github.com/gin-gonic/gin"
)

func accessKeyMiddleware(handler *action.Handler, now func() time.Time, server *fdsnwsgo.FDSNWS, services ...fdsnwsgo.Service) gin.HandlerFunc {
	queryModules := make(map[string]fdsnwsgo.Module, len(services))
	for _, service := range services {
		queryPath := fmt.Sprintf("%s/%s/%s/query", basePath, service.Name(), service.MajorVersion())
		queryModules[queryPath] = service.Name()
	}
	return func(ctx *gin.Context) {
		reqPath := strings.TrimSuffix(ctx.Request.URL.Path, "/")
		module, protected := queryModules[reqPath]
		if !protected {
			ctx.Next()
			return
		}

		submittedAt := now()
		status := http.StatusUnauthorized
		detail := "FDSN access key is required"
		fields := strings.Fields(ctx.GetHeader("Authorization"))
		if len(fields) == 2 && strings.EqualFold(fields[0], "Bearer") {
			_, err := handler.FDSNAccessKeyAuthenticate(ctx.Request.Context(), fields[1], ctx.ClientIP(), ctx.Request.UserAgent(), submittedAt)
			switch {
			case err == nil:
				ctx.Next()
				return
			case errors.Is(err, action.ErrInvalidFDSNAccessKey):
				detail = "Invalid, expired or revoked FDSN access key"
			default:
				_ = ctx.Error(err)
				status = http.StatusServiceUnavailable
				detail = "FDSN access validation is unavailable"
			}
		}

		response := server.ErrorResponse(fdsnwsgo.Request{
			URL: ctx.Request.URL.String(), ReceivedAt: submittedAt,
		}, module, status, http.StatusText(status), detail)
		if status == http.StatusUnauthorized {
			response.Headers = map[string]string{"WWW-Authenticate": bearerChallenge}
		}
		ctx.Abort()
		if err := fdsnwsgo.WriteHTTPResponse(ctx.Writer, response); err != nil {
			_ = ctx.Error(err)
		}
	}
}
