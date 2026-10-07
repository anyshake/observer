package notifications

import (
	"bufio"
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/notification"
	"github.com/gin-gonic/gin"
)

func TestNotificationEventAndClientCancel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	hub := notification.NewHub(1)
	shutdownCtx, cancelShutdown := context.WithCancel(context.Background())
	defer cancelShutdown()

	engine := gin.New()
	Setup(engine.Group("/api"), hub, func(ctx *gin.Context) { ctx.Next() }, shutdownCtx)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := &http.Server{Handler: engine}
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(func() { _ = server.Close() })

	requestCtx, cancelRequest := context.WithCancel(context.Background())
	request, err := http.NewRequestWithContext(requestCtx, http.MethodGet, "http://"+listener.Addr().String()+"/api/notifications", nil)
	if err != nil {
		t.Fatal(err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	reader := bufio.NewReader(response.Body)
	if line, err := reader.ReadString('\n'); err != nil || line != ": connected\n" {
		t.Fatalf("connected line = %q, %v", line, err)
	}
	hub.Publish(notification.Event{ID: "evt-1", ServiceID: "quakesense", Message: "triggered", Level: notification.LevelWarning, OccurredAt: 10})
	var payload string
	deadline := time.Now().Add(time.Second)
	for !strings.Contains(payload, "triggered") {
		if time.Now().After(deadline) {
			t.Fatalf("event was not delivered: %q", payload)
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			t.Fatal(err)
		}
		payload += line
	}
	cancelRequest()
}

func TestShutdownWithConnectedClient(t *testing.T) {
	gin.SetMode(gin.TestMode)

	shutdownCtx, cancelShutdown := context.WithCancel(context.Background())
	defer cancelShutdown()

	engine := gin.New()
	Setup(engine.Group("/api"), notification.NewHub(1), func(ctx *gin.Context) {
		ctx.Next()
	}, shutdownCtx)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}

	server := &http.Server{Handler: engine}
	serveErr := make(chan error, 1)
	go func() {
		serveErr <- server.Serve(listener)
	}()
	t.Cleanup(func() {
		_ = server.Close()
	})

	response, err := http.Get("http://" + listener.Addr().String() + "/api/notifications")
	if err != nil {
		t.Fatalf("failed to connect to SSE endpoint: %v", err)
	}
	defer response.Body.Close()

	connectedLine, err := bufio.NewReader(response.Body).ReadString('\n')
	if err != nil {
		t.Fatalf("failed to read SSE connection event: %v", err)
	}
	if connectedLine != ": connected\n" {
		t.Fatalf("unexpected SSE connection event: %q", connectedLine)
	}

	cancelShutdown()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		t.Fatalf("failed to shut down server with connected SSE client: %v", err)
	}

	select {
	case err := <-serveErr:
		if !errors.Is(err, http.ErrServerClosed) {
			t.Fatalf("unexpected server error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("server did not stop after shutdown")
	}
}
