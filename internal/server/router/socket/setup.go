package socket

import (
	"net/http"
	"time"

	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/server/response"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/message"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/samber/lo"
)

func Setup(routerGroup *gin.RouterGroup, timeSource *timesource.Source, hardware hardware.IHardware, jwtMiddleware gin.HandlerFunc) {
	s := socket{
		messageBus:    message.NewBus[explorer.EventHandler](LOG_PREFIX, 65535),
		historyBuffer: make([]buffer, 0, HISTORY_BUFFER_SIZE),
	}
	hardware.Subscribe(LOG_PREFIX, func(t time.Time, di *explorer.DeviceConfig, dv *explorer.DeviceVariable, cd []explorer.ChannelData) {
		s.messageBus.Publish(t, di, dv, cd)
		s.storeHistory(t, di, cd)
	})

	routerGroup.GET("/socket", jwtMiddleware, func(ctx *gin.Context) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
			EnableCompression: true,
			Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
				logger.GetLogger(LOG_PREFIX).Errorf("websocket error %d: %s", status, reason)
				response.Error(ctx, status, "websocket error")
			},
			CheckOrigin: func(r *http.Request) bool { return true },
		}

		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			logger.GetLogger(LOG_PREFIX).Errorf("failed to upgrade connection: %v", err)
			return
		}
		defer conn.Close()

		s.handleWebSocket(ctx, conn, timeSource)
	})
}

func (s *socket) storeHistory(t time.Time, di *explorer.DeviceConfig, cd []explorer.ChannelData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.historyBuffer) >= HISTORY_BUFFER_SIZE {
		s.historyBuffer = s.historyBuffer[1:]
	}
	s.historyBuffer = append(s.historyBuffer, buffer{
		Timestamp:   t.UnixMilli(),
		SampleRate:  di.GetSampleRate(),
		ChannelData: cd,
	})
}

func (s *socket) sendHistory(conn *websocket.Conn, timeSource *timesource.Source) error {
	s.mu.Lock()
	historyMessages := lo.Map(s.historyBuffer, func(history buffer, _ int) map[string]any {
		return map[string]any{
			"current_time": timeSource.Now().UnixMilli(),
			"record_time":  history.Timestamp,
			"sample_rate":  history.SampleRate,
			"channel_data": lo.SliceToMap(history.ChannelData, func(v explorer.ChannelData) (string, any) {
				return v.ChannelCode, map[string]any{
					"channel_id":   v.ChannelId,
					"channel_code": v.ChannelCode,
					"data_type":    v.DataType,
					"data_array":   v.Data,
				}
			}),
		}
	})
	s.mu.Unlock()

	for _, message := range historyMessages {
		s.mu.Lock()
		_ = conn.WriteJSON(message)
		s.mu.Unlock()
	}
	return nil
}

func (s *socket) handleWebSocket(_ *gin.Context, conn *websocket.Conn, timeSource *timesource.Source) {
	clientID := conn.RemoteAddr().String()
	subscribedAt := time.Now()
	logger.GetLogger(LOG_PREFIX).Infof("%s - client subscribed to message bus", clientID)

	callbackFn := func(t time.Time, di *explorer.DeviceConfig, dv *explorer.DeviceVariable, cd []explorer.ChannelData) {
		data := map[string]any{
			"current_time": timeSource.Now().UnixMilli(),
			"record_time":  t.UnixMilli(),
			"sample_rate":  di.GetSampleRate(),
			"channel_data": lo.SliceToMap(cd, func(v explorer.ChannelData) (string, any) {
				return v.ChannelCode, map[string]any{
					"channel_id":   v.ChannelId,
					"channel_code": v.ChannelCode,
					"data_type":    v.DataType,
					"data_array":   v.Data,
				}
			}),
		}
		s.mu.Lock()
		_ = conn.WriteJSON(data)
		s.mu.Unlock()
	}

	if err := s.messageBus.Subscribe(clientID, callbackFn); err != nil {
		logger.GetLogger(LOG_PREFIX).Errorf("failed to subscribe: %v", err)
		return
	}
	defer s.messageBus.Unsubscribe(clientID)

	for {
		_, dataBytes, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if string(dataBytes) == "client hello" {
			if err := s.sendHistory(conn, timeSource); err != nil {
				break
			}
		}
	}

	duration := time.Since(subscribedAt).Seconds()
	logger.GetLogger(LOG_PREFIX).Infof("%s - unsubscribed after %f seconds", clientID, duration)
}
