package socket

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/hardware/explorer"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/metadata"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TestWebSocketStream(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Init()
	hardware := &fakeHardware{}
	router := gin.New()
	Setup(router.Group("/"), timesource.New(func() time.Time {
		return time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	}), hardware, func(ctx *gin.Context) {
		if ctx.GetHeader("Authorization") != "Bearer good" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	})

	server := httptest.NewServer(router)
	t.Cleanup(server.Close)
	endpoint := "ws" + strings.TrimPrefix(server.URL, "http") + "/socket"

	plain := httptest.NewRecorder()
	router.ServeHTTP(plain, httptest.NewRequest(http.MethodGet, "/socket", nil))
	if plain.Code == 0 {
		t.Fatal("plain request produced no status")
	}

	reject, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := reject.WriteMessage(websocket.TextMessage, []byte("bad")); err != nil {
		t.Fatal(err)
	}
	_ = reject.SetReadDeadline(time.Now().Add(time.Second))
	if _, _, err := reject.ReadMessage(); err == nil {
		t.Fatal("rejected client stayed open")
	}
	_ = reject.Close()

	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if err := conn.WriteMessage(websocket.TextMessage, []byte("good")); err != nil {
		t.Fatal(err)
	}
	event := explorer.Event{
		Timestamp:  time.UnixMilli(1_700_000_000_000),
		SampleRate: 100,
		ChannelData: []explorer.ChannelData{{
			ChannelCode: "EHZ",
			ChannelId:   1,
			DataType:    "int32",
			Data:        []int32{4, 5},
		}},
	}
	deadline := time.Now().Add(time.Second)
	for hardware.handler == nil {
		if time.Now().After(deadline) {
			t.Fatal("hardware subscription was not registered")
		}
		time.Sleep(5 * time.Millisecond)
	}
	hardware.handler(event)
	if err := conn.WriteMessage(websocket.TextMessage, []byte("client hello")); err != nil {
		t.Fatal(err)
	}
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, history, err := conn.ReadMessage()
	if err != nil || !strings.Contains(string(history), "EHZ") {
		t.Fatalf("history = %s, %v", history, err)
	}
	hardware.handler(event)
	_, live, err := conn.ReadMessage()
	if err != nil || !strings.Contains(string(live), "EHZ") {
		t.Fatalf("live = %s, %v", live, err)
	}

	for range HISTORY_BUFFER_SIZE + 1 {
		hardware.handler(event)
	}
}

func TestWebSocketSubscribeFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Init()
	router := gin.New()
	Setup(router.Group("/"), timesource.New(time.Now), &fakeHardware{subscribeErr: context.Canceled}, func(ctx *gin.Context) {})
	server := httptest.NewServer(router)
	t.Cleanup(server.Close)
	conn, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"/socket", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if err := conn.WriteMessage(websocket.TextMessage, []byte("good")); err != nil {
		t.Fatal(err)
	}
	time.Sleep(20 * time.Millisecond)
}

type fakeHardware struct {
	handler      explorer.EventHandler
	subscribeErr error
}

func (f *fakeHardware) Open(context.Context) (context.Context, context.CancelFunc, error) {
	return context.Background(), func() {}, nil
}
func (fakeHardware) Close() error { return nil }
func (fakeHardware) Flush() error { return nil }
func (f *fakeHardware) Subscribe(_ string, handler explorer.EventHandler) error {
	if f.subscribeErr != nil {
		return f.subscribeErr
	}
	f.handler = handler
	return nil
}
func (fakeHardware) Unsubscribe(string) error { return nil }
func (fakeHardware) SubscribeRealtime(string, explorer.EventHandler) error {
	return nil
}
func (fakeHardware) UnsubscribeRealtime(string) error { return nil }
func (fakeHardware) GetConfig() explorer.DeviceConfig { return explorer.DeviceConfig{} }
func (fakeHardware) GetStatus() explorer.DeviceStatus { return explorer.DeviceStatus{} }
func (fakeHardware) GetCoordinates(bool) (float64, float64, float64, error) {
	return 0, 0, 0, nil
}
func (fakeHardware) GetTemperature() (float64, error) { return 0, nil }
func (fakeHardware) GetDeviceId() string              { return "device" }
func (fakeHardware) GetMetadata(string, string, string, string, string, string, string, bool) (*metadata.Render, error) {
	return &metadata.Render{}, nil
}

func TestNoopResponseWriter(t *testing.T) {
	writer := &noopResponseWriter{header: make(http.Header)}
	writer.Header().Set("X-Test", "1")
	n, err := writer.Write([]byte("ab"))
	if n != 2 || err != nil {
		t.Fatal(err)
	}
	writer.WriteHeader(http.StatusNoContent)
	if writer.code != http.StatusNoContent || writer.Header().Get("X-Test") != "1" {
		t.Fatal("writer did not record the response")
	}
}
