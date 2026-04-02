package auth_jwt

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewWebsocketAuthAdapter() gin.HandlerFunc {
	return func(c *gin.Context) {
		upgrade := strings.ToLower(c.GetHeader("Upgrade"))
		connection := strings.ToLower(c.GetHeader("Connection"))

		if upgrade != "websocket" || !strings.Contains(connection, "upgrade") {
			c.Next()
			return
		}

		protocol := c.GetHeader("Sec-WebSocket-Protocol")
		if protocol == "" {
			c.Next()
			return
		}

		if token := extractBearerToken(protocol); token != "" {
			c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		}

		c.Next()
	}
}

func extractBearerToken(protocols string) string {
	for _, p := range strings.Split(protocols, ",") {
		p = strings.TrimSpace(p)
		if strings.HasPrefix(p, "Bearer#") {
			return strings.TrimPrefix(p, "Bearer#")
		}
	}

	return ""
}
