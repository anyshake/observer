package socket

import (
	"encoding/json"
	"net/http"

	"github.com/bclswl0827/observer/app"
	"github.com/bclswl0827/observer/publisher"
	"github.com/bclswl0827/observer/server/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func (s *Socket) RegisterModule(rg *gin.RouterGroup, options *app.ServerOptions) {
	rg.GET("/socket", func(c *gin.Context) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize: 1024, WriteBufferSize: 1024, EnableCompression: true,
			Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
				response.Error(c, http.StatusBadRequest)
			},
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		// Properly close connection
		go func() {
			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					break
				}
			}

			conn.Close()
		}()

		// Write when new message arrived
		publisher.Subscribe(
			&options.FeatureOptions.Status.Geophone,
			func(gp *publisher.Geophone) error {
				data, err := json.Marshal(gp)
				if err != nil {
					conn.Close()
					return err
				}

				err = conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					conn.Close()
					return err
				}

				return nil
			},
		)
	})
}
