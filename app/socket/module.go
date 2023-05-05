package socket

import (
	"encoding/json"
	"net/http"
	"time"

	"com.geophone.observer/app"
	"com.geophone.observer/server/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var lastTimestamp int64

func (s *Socket) RegisterModule(rg *gin.RouterGroup, options *app.ServerOptions) {
	rg.GET("/socket", func(c *gin.Context) {
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
			if lastTimestamp == options.Message.Acceleration.Timestamp {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			data, err := json.Marshal(options.Message)
			if err != nil {
				return
			}

			err = conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				return
			}

			lastTimestamp = options.Message.Acceleration.Timestamp
		}
	})
}
