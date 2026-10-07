package socket

import (
	"net/http"
	"sync"
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
		messageBus:     message.NewBus[explorer.Event](LOG_PREFIX),
		historyBuffer:  make([]buffer, HISTORY_BUFFER_SIZE),
		tokenValidator: newTokenValidator(jwtMiddleware),
	}
	if err := hardware.Subscribe(LOG_PREFIX, func(event explorer.Event) {
		s.messageBus.Publish(event)
		s.storeHistory(event)
	}); err != nil {
		logger.GetLogger(LOG_PREFIX).Errorf("failed to subscribe to hardware message bus: %v", err)
	}

	routerGroup.GET("/socket", func(ctx *gin.Context) {
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

func (s *socket) storeHistory(event explorer.Event) {
	channelData := make([]explorer.ChannelData, len(event.ChannelData))
	for i := range event.ChannelData {
		channelData[i] = event.ChannelData[i]
		channelData[i].Data = append([]int32(nil), event.ChannelData[i].Data...)
	}

	s.historyMu.Lock()
	s.historyBuffer[s.historyPos] = buffer{
		Timestamp:   event.Timestamp.UnixMilli(),
		SampleRate:  event.SampleRate,
		ChannelData: channelData,
	}
	s.historyPos++
	if s.historyPos == HISTORY_BUFFER_SIZE {
		s.historyPos = 0
	}
	if s.historyLen < HISTORY_BUFFER_SIZE {
		s.historyLen++
	}
	s.historyMu.Unlock()
}

func (s *socket) historyAt(index int) buffer {
	start := 0
	if s.historyLen == HISTORY_BUFFER_SIZE {
		start = s.historyPos
	}
	return s.historyBuffer[(start+index)%HISTORY_BUFFER_SIZE]
}

func (s *socket) sendHistory(conn *websocket.Conn, writeMu *sync.Mutex, timeSource *timesource.Source) error {
	s.historyMu.RLock()
	historyMessages := make([]map[string]any, s.historyLen)
	for i := range historyMessages {
		history := s.historyAt(i)
		historyMessages[i] = map[string]any{
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
	}
	s.historyMu.RUnlock()

	for _, message := range historyMessages {
		if err := writeWebSocketJSON(conn, writeMu, message); err != nil {
			return err
		}
	}
	return nil
}

func writeWebSocketJSON(conn *websocket.Conn, writeMu *sync.Mutex, data any) error {
	writeMu.Lock()
	defer writeMu.Unlock()

	if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
		_ = conn.Close()
		return err
	}
	if err := conn.WriteJSON(data); err != nil {
		_ = conn.Close()
		return err
	}
	return nil
}

func (s *socket) handleWebSocket(_ *gin.Context, conn *websocket.Conn, timeSource *timesource.Source) {
	clientID := conn.RemoteAddr().String()
	logger.GetLogger(LOG_PREFIX).Infof("%s - client connected, waiting for authentication", clientID)

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	_, tokenBytes, err := conn.ReadMessage()
	conn.SetReadDeadline(time.Time{})
	if err != nil {
		logger.GetLogger(LOG_PREFIX).Warnf("%s - authentication timeout or read error: %v", clientID, err)
		return
	}
	if !s.tokenValidator(string(tokenBytes)) {
		logger.GetLogger(LOG_PREFIX).Warnf("%s - authentication failed", clientID)
		return
	}

	subscribedAt := time.Now()
	logger.GetLogger(LOG_PREFIX).Infof("%s - authenticated and subscribed to message bus", clientID)

	var writeMu sync.Mutex
	callbackFn := func(event explorer.Event) {
		data := map[string]any{
			"current_time": timeSource.Now().UnixMilli(),
			"record_time":  event.Timestamp.UnixMilli(),
			"sample_rate":  event.SampleRate,
			"channel_data": lo.SliceToMap(event.ChannelData, func(v explorer.ChannelData) (string, any) {
				return v.ChannelCode, map[string]any{
					"channel_id":   v.ChannelId,
					"channel_code": v.ChannelCode,
					"data_type":    v.DataType,
					"data_array":   v.Data,
				}
			}),
		}
		if err := writeWebSocketJSON(conn, &writeMu, data); err != nil {
			logger.GetLogger(LOG_PREFIX).Warnf("%s - failed to write websocket data: %v", clientID, err)
		}
	}

	if err := s.messageBus.Subscribe(clientID, message.SubscriptionOptions{
		BufferSize: 2,
		Overflow:   message.OverflowDropOldest,
		OnError: func(err error) {
			logger.GetLogger(LOG_PREFIX).Warnf("%s - websocket subscription closed: %v", clientID, err)
			_ = conn.Close()
		},
	}, callbackFn); err != nil {
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
			if err := s.sendHistory(conn, &writeMu, timeSource); err != nil {
				logger.GetLogger(LOG_PREFIX).Warnf("%s - failed to write websocket history: %v", clientID, err)
				break
			}
		}
	}

	duration := time.Since(subscribedAt).Seconds()
	logger.GetLogger(LOG_PREFIX).Infof("%s - unsubscribed after %f seconds", clientID, duration)
}
