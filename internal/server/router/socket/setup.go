package socket

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/anyshake/observer/internal/hardware"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/internal/server/response"
	"github.com/anyshake/observer/pkg/ctxutil"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/message"
	"github.com/anyshake/observer/pkg/ringbuf"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/samber/lo"
)

// Clients only send a JWT or short commands, including when compression is negotiated.
const maxClientMessageSize = 8 * 1024

func Setup(routerGroup *gin.RouterGroup, timeSource *timesource.Source, hardware hardware.IHardware, jwtMiddleware gin.HandlerFunc, shutdownCtx context.Context) func(context.Context) error {
	s := socket{
		messageBus:     message.NewBus[explorer.Event](LOG_PREFIX),
		history:        ringbuf.New[buffer](HISTORY_BUFFER_SIZE),
		tokenValidator: newTokenValidator(jwtMiddleware),
	}
	subscribeErr := hardware.Subscribe(LOG_PREFIX, func(event explorer.Event) {
		s.messageBus.Publish(event)
		s.storeHistory(event)
	})
	if subscribeErr != nil {
		logger.GetLogger(LOG_PREFIX).Errorf("failed to subscribe to hardware message bus: %v", subscribeErr)
	}

	routerGroup.GET("/socket", func(ctx *gin.Context) {
		s.mu.Lock()
		if s.closed {
			s.mu.Unlock()
			response.Error(ctx, http.StatusServiceUnavailable, "server is shutting down")
			return
		}
		s.handlers.Add(1)
		s.mu.Unlock()
		defer s.handlers.Done()

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
		stopClose := ctxutil.AfterFunc(shutdownCtx, func() { _ = conn.Close() })
		defer stopClose()

		s.handleWebSocket(ctx, conn, timeSource)
	})

	return func(ctx context.Context) error {
		s.mu.Lock()
		wasClosed := s.closed
		s.closed = true
		s.mu.Unlock()
		if subscribeErr == nil && !wasClosed {
			_ = hardware.Unsubscribe(LOG_PREFIX)
		}

		finished := make(chan struct{})
		go func() {
			s.handlers.Wait()
			close(finished)
		}()
		select {
		case <-finished:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *socket) storeHistory(event explorer.Event) {
	channelData := make([]explorer.ChannelData, len(event.ChannelData))
	for i := range event.ChannelData {
		channelData[i] = event.ChannelData[i]
		channelData[i].Data = append([]int32(nil), event.ChannelData[i].Data...)
	}

	s.history.Push(buffer{
		Timestamp:   event.Timestamp.UnixMilli(),
		SampleRate:  event.SampleRate,
		ChannelData: channelData,
	})
}

func (s *socket) sendHistory(conn *websocket.Conn, writeMu *sync.Mutex, timeSource *timesource.Source) error {
	historyValues := s.history.Values()
	historyMessages := make([]map[string]any, len(historyValues))
	for i, history := range historyValues {
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

func readWebSocketMessage(conn *websocket.Conn) ([]byte, error) {
	_, reader, err := conn.NextReader()
	if err != nil {
		return nil, err
	}
	// SetReadLimit bounds wire payloads; this also bounds decompressed messages.
	data, err := io.ReadAll(io.LimitReader(reader, maxClientMessageSize+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxClientMessageSize {
		_ = conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseMessageTooBig, "message too large"),
			time.Now().Add(time.Second))
		return nil, websocket.ErrReadLimit
	}
	return data, nil
}

func (s *socket) handleWebSocket(_ *gin.Context, conn *websocket.Conn, timeSource *timesource.Source) {
	clientID := conn.RemoteAddr().String()
	logger.GetLogger(LOG_PREFIX).Infof("%s - client connected, waiting for authentication", clientID)

	conn.SetReadLimit(maxClientMessageSize)
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	tokenBytes, err := readWebSocketMessage(conn)
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
		dataBytes, err := readWebSocketMessage(conn)
		if err != nil {
			break
		}
		if string(dataBytes) == "ping" {
			if err := writeWebSocketJSON(conn, &writeMu, "pong"); err != nil {
				logger.GetLogger(LOG_PREFIX).Warnf("%s - failed to write websocket heartbeat: %v", clientID, err)
				break
			}
			continue
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
