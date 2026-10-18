package seedlink

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/anyshake/observer/pkg/logger"
	"github.com/bclswl0827/slgo/handlers"
)

func TestSeedLinkWriteRefreshesExpiredDeadline(t *testing.T) {
	logger.Init()
	server, peer := net.Pipe()
	defer server.Close()
	defer peer.Close()
	client := &handlers.SeedLinkClient{Conn: server}
	hooks := &hooks{clients: newSeedLinkClientRegistry()}
	hooks.OnConnection(client)
	defer hooks.OnClose(client)
	// Model a deadline that expired while a finite history query was running;
	// no live callback or previous successful write can refresh it.
	if err := server.SetWriteDeadline(time.Now().Add(-time.Second)); err != nil {
		t.Fatal(err)
	}
	if err := peer.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() {
		_, err := client.Write([]byte("END"))
		done <- err
	}()
	output := make([]byte, 3)
	if _, err := io.ReadFull(peer, output); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if string(output) != "END" {
		t.Fatalf("reply = %q, want END", output)
	}
}
