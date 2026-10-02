package notifications

import (
	"context"
	"io"
	"time"

	"github.com/anyshake/observer/internal/notification"
	"github.com/gin-gonic/gin"
)

const heartbeatInterval = 20 * time.Second

func Setup(routerGroup *gin.RouterGroup, hub *notification.Hub, jwtMiddleware gin.HandlerFunc, shutdownCtx context.Context) {
	routerGroup.GET("/notifications", jwtMiddleware, func(ctx *gin.Context) {
		events, unsubscribe := hub.Subscribe()
		defer unsubscribe()

		ctx.Header("Content-Type", "text/event-stream")
		ctx.Header("Cache-Control", "no-cache")
		ctx.Header("Connection", "keep-alive")
		ctx.Header("X-Accel-Buffering", "no")

		_, _ = io.WriteString(ctx.Writer, ": connected\n\n")
		ctx.Writer.Flush()

		heartbeat := time.NewTicker(heartbeatInterval)
		defer heartbeat.Stop()

		ctx.Stream(func(writer io.Writer) bool {
			select {
			case event, ok := <-events:
				if !ok {
					return false
				}
				ctx.SSEvent("service-notification", event)
			case <-heartbeat.C:
				if _, err := io.WriteString(writer, ": keepalive\n\n"); err != nil {
					return false
				}
			case <-ctx.Request.Context().Done():
				return false
			case <-shutdownCtx.Done():
				return false
			}

			return true
		})
	})
}
