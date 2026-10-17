package dnsquery

import (
	"context"
	"time"

	"github.com/miekg/dns"
)

// Closing the connection on cancellation also interrupts reads when ctx has no
// deadline; ExchangeContext by itself only propagates deadlines to socket I/O.
func exchangeContext(ctx context.Context, client *dns.Client, msg *dns.Msg, server string) (*dns.Msg, error) {
	conn, err := client.DialContext(ctx, server)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	stop := closeOnCancel(ctx, func() { _ = conn.Close() })
	defer stop()
	response, _, err := client.ExchangeWithConnContext(ctx, msg, conn)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return response, err
}

func closeOnCancel(ctx context.Context, closeConn func()) func() {
	done := make(chan struct{})
	stop := context.AfterFunc(ctx, func() {
		closeConn()
		close(done)
	})
	return func() {
		if !stop() {
			<-done
		}
	}
}

func queryContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return context.WithTimeout(ctx, timeout)
}
