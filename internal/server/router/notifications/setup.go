package notifications

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/anyshake/observer/internal/notification"
	"github.com/anyshake/observer/internal/server/middleware/auth_jwt"
	"github.com/gin-gonic/gin"
)

const (
	heartbeatInterval = 20 * time.Second
	writeTimeout      = 5 * time.Second
)

func Setup(routerGroup *gin.RouterGroup, hub *notification.Hub, jwtMiddleware gin.HandlerFunc, shutdownCtx context.Context) {
	routerGroup.GET("/notifications", jwtMiddleware, func(ctx *gin.Context) {
		events, unsubscribe := hub.Subscribe()
		defer unsubscribe()

		ctx.Header("Content-Type", "text/event-stream")
		ctx.Header("Cache-Control", "no-cache")
		ctx.Header("Connection", "keep-alive")
		ctx.Header("X-Accel-Buffering", "no")

		controller := http.NewResponseController(ctx.Writer)
		setWriteDeadline := func() bool {
			err := controller.SetWriteDeadline(time.Now().Add(writeTimeout))
			return err == nil || errors.Is(err, http.ErrNotSupported)
		}
		flush := func() bool {
			if err := controller.Flush(); err != nil {
				return false
			}
			// Only an actual write is timed; idle streams wait for the heartbeat.
			err := controller.SetWriteDeadline(time.Time{})
			return err == nil || errors.Is(err, http.ErrNotSupported)
		}
		if !setWriteDeadline() {
			return
		}
		if _, err := io.WriteString(ctx.Writer, ": connected\n\n"); err != nil {
			return
		}
		if !flush() {
			return
		}

		heartbeat := time.NewTicker(heartbeatInterval)
		defer heartbeat.Stop()

		ctx.Stream(func(writer io.Writer) bool {
			select {
			case event, ok := <-events:
				if !ok {
					return false
				}
				if event.AdminOnly && !ctx.GetBool(auth_jwt.IsAdminKey) {
					return true
				}
				if !setWriteDeadline() {
					return false
				}
				ctx.SSEvent("service-notification", event)
				if ctx.IsAborted() {
					return false
				}
			case <-heartbeat.C:
				if !setWriteDeadline() {
					return false
				}
				if _, err := io.WriteString(writer, ": keepalive\n\n"); err != nil {
					return false
				}
			case <-ctx.Request.Context().Done():
				return false
			case <-shutdownCtx.Done():
				return false
			}

			return flush()
		})
	})
}
