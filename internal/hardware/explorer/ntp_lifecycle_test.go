package explorer

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"

	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/ntpclient"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/anyshake/observer/pkg/transport"
	"github.com/rs/zerolog"
)

type ntpLifecycleTransport struct {
	transport.ITransport
	closed      int
	beforeClose func()
}

func (*ntpLifecycleTransport) Open() error  { return nil }
func (*ntpLifecycleTransport) Flush() error { return nil }
func (p *ntpLifecycleTransport) Close() error {
	if p.beforeClose != nil {
		p.beforeClose()
	}
	p.closed++
	return nil
}

func TestProtocolCloseStopsNTPLifecycle(t *testing.T) {
	for _, version := range []string{"v1", "v2", "v3"} {
		t.Run(version, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				client, err := ntpclient.New([]string{"ntp://time.example.test"}, 0, 1, nil)
				if err != nil {
					t.Fatal(err)
				}
				defer client.Close()
				done := make(chan struct{})
				transport := &ntpLifecycleTransport{beforeClose: func() {
					if ctx.Err() == nil {
						t.Error("transport closed before canceling synchronization")
					}
					if _, _, err := client.QueryContext(ctx); !errors.Is(err, ntpclient.ErrClosed) {
						t.Errorf("NTP client was not closed: %v", err)
					}
					select {
					case <-done:
					default:
						t.Error("transport closed before synchronization worker exited")
					}
				}}
				var closeProtocol func() error
				switch version {
				case "v1":
					protocol := &ExplorerProtoImplV1{Transport: transport, ntpClient: client, cancelFn: cancel, ntpDone: done}
					closeProtocol = protocol.Close
				case "v2":
					protocol := &ExplorerProtoImplV2{Transport: transport, ntpClient: client, cancelFn: cancel, ntpDone: done}
					closeProtocol = protocol.Close
				case "v3":
					protocol := &ExplorerProtoImplV3{Transport: transport, ntpClient: client, cancelFn: cancel, ntpDone: done}
					closeProtocol = protocol.Close
				}
				go func() {
					<-ctx.Done()
					time.Sleep(time.Second)
					close(done)
				}()
				if err := closeProtocol(); err != nil {
					t.Fatal(err)
				}
				if transport.closed != 1 {
					t.Fatalf("transport closed %d times, want 1", transport.closed)
				}
			})
		})
	}
}

func TestProtocolCanceledStartupCleansUp(t *testing.T) {
	for _, version := range []string{"v1", "v2", "v3"} {
		t.Run(version, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				transport := &ntpLifecycleTransport{}
				log := &logger.Adapter{Logger: zerolog.Nop()}
				base := time.Unix(1000, 123)
				source := timesource.New(func() time.Time { return base })
				options := NtpOptions{Pool: []string{"ntp://time.example.test"}, ReadTimeout: 1}
				var open func(context.Context) (context.Context, context.CancelFunc, error)
				var getClient func() *ntpclient.Client
				switch version {
				case "v1":
					protocol := &ExplorerProtoImplV1{Transport: transport, Logger: log, TimeSource: source, NtpOptions: options}
					open, getClient = protocol.Open, func() *ntpclient.Client { return protocol.ntpClient }
				case "v2":
					protocol := &ExplorerProtoImplV2{Transport: transport, Logger: log, TimeSource: source, NtpOptions: options}
					open, getClient = protocol.Open, func() *ntpclient.Client { return protocol.ntpClient }
				case "v3":
					protocol := &ExplorerProtoImplV3{Transport: transport, Logger: log, TimeSource: source, NtpOptions: options}
					open, getClient = protocol.Open, func() *ntpclient.Client { return protocol.ntpClient }
				}
				if _, _, err := open(ctx); !errors.Is(err, context.Canceled) {
					t.Fatalf("canceled Open error = %v", err)
				}
				if transport.closed != 1 {
					t.Fatalf("failed startup closed transport %d times, want 1", transport.closed)
				}
				if client := getClient(); client == nil {
					t.Fatal("startup did not create an NTP client")
				} else if _, _, err := client.QueryContext(ctx); !errors.Is(err, ntpclient.ErrClosed) {
					t.Fatalf("failed startup left NTP client open: %v", err)
				}
				if got := source.Now(); !got.Equal(base) {
					t.Errorf("canceled startup changed the time source: %v", got)
				}
			})
		})
	}
}
