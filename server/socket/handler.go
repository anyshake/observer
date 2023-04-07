package socket

import (
	"encoding/json"
	"net/http"
	"time"

	"com.geophone.observer/features/collector"
	"com.geophone.observer/server/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func WebsocketHandler(c *gin.Context, m *collector.Message) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize: 1024, WriteBufferSize: 1024, EnableCompression: true,
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			response.ErrorHandler(c, http.StatusBadRequest)
		},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	for {
		data, _ := json.Marshal(m)
		conn.WriteMessage(websocket.TextMessage, data)
		time.Sleep(10 * time.Second)
	}
}
