package socket

import (
	"encoding/json"
	"net/http"

	"github.com/alphadose/haxmap"
	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/server/response"
	"github.com/anyshake/observer/services"
	"github.com/anyshake/observer/utils/logger"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	messagebus "github.com/vardius/message-bus"
)

func (s *Socket) Bind(rg *gin.RouterGroup, jwtHandler *jwt.GinJWTMiddleware, options *services.Options) error {
	s.subscribers = haxmap.New[string, explorer.ExplorerEventHandler]()
	s.messageBus = messagebus.New(65535)

	// Forward events to internal message bus
	var explorerDeps *explorer.ExplorerDependency
	err := options.Dependency.Invoke(func(deps *explorer.ExplorerDependency) error {
		explorerDeps = deps
		return nil
	})
	if err != nil {
		logger.GetLogger(s.GetApiName()).Errorln(err)
		return err
	}
	explorerDriver := explorer.ExplorerDriver(&explorer.ExplorerDriverImpl{})
	explorerDriver.Subscribe(
		explorerDeps,
		s.GetApiName(),
		func(data *explorer.ExplorerData) {
			s.messageBus.Publish(s.GetApiName(), data)
			s.historyBuffer[s.historyBufferIndex] = *data
			s.historyBufferIndex = (s.historyBufferIndex + 1) % HISTORY_BUFFER_SIZE
		},
	)

	var handlerFunc []gin.HandlerFunc
	if options.Config.Server.Restrict {
		handlerFunc = append(handlerFunc, jwtHandler.MiddlewareFunc())
	}
	handlerFunc = append(handlerFunc, func(c *gin.Context) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize: 1024, WriteBufferSize: 1024, EnableCompression: true,
			Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
				logger.GetLogger(s.GetApiName()).Errorf("websocket error, code %d, %s", status, reason)
				response.Message(c, options.TimeSource, "websocket error occurred", http.StatusBadRequest, nil)
			},
			CheckOrigin: func(r *http.Request) bool { return true },
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.GetLogger(s.GetApiName()).Errorln(err)
			return
		}
		defer conn.Close()

		// Subscribe to the internal message bus
		clienrId := conn.RemoteAddr().String()
		handler := func(data *explorer.ExplorerData) {
			dataBytes, err := json.Marshal(data)
			if err != nil {
				logger.GetLogger(s.GetApiName()).Errorln(err)
				return
			}
			err = conn.WriteMessage(websocket.TextMessage, dataBytes)
			if err != nil {
				logger.GetLogger(s.GetApiName()).Errorln(err)
				return
			}
		}
		err = s.subscribe(clienrId, handler)
		if err != nil {
			logger.GetLogger(s.GetApiName()).Errorln(err)
			return
		}
		defer s.unsubscribe(clienrId)

		// Listen for incoming messages
		for {
			_, dataBytes, err := conn.ReadMessage()
			if err != nil {
				return
			}
			// Respond with history buffer if the client sends a "client hello" message
			if string(dataBytes) == "client hello" {
				for _, buffer := range s.historyBuffer {
					if buffer.Timestamp == 0 {
						continue
					}
					dataBytes, err := json.Marshal(buffer)
					if err != nil {
						logger.GetLogger(s.GetApiName()).Errorln(err)
						return
					}
					err = conn.WriteMessage(websocket.TextMessage, dataBytes)
					if err != nil {
						logger.GetLogger(s.GetApiName()).Errorln(err)
						return
					}
				}
			}
		}
	})

	rg.GET("/socket", handlerFunc...)
	return nil
}
