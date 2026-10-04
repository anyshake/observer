package ntpclient

import (
	"context"
	"net"
	"net/netip"
	"time"

	"github.com/beevik/ntp"
)

// Use the library's packet handling, with cancellation covering DNS and UDP I/O.
func queryWithContext(ctx context.Context, server string, options ntp.QueryOptions, query queryFunc) (*ntp.Response, error) {
	timeout := options.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second // Match ntp.QueryOptions' default.
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	dial := options.Dialer
	if dial == nil {
		dial = func(local, remote string) (net.Conn, error) {
			dialer := net.Dialer{}
			if local != "" {
				address, err := netip.ParseAddr(local)
				if err != nil {
					return nil, err
				}
				dialer.LocalAddr = net.UDPAddrFromAddrPort(netip.AddrPortFrom(address, 0))
			}
			return dialer.DialContext(ctx, "udp", remote)
		}
	}

	var stop func() bool
	done := make(chan struct{})
	options.Dialer = func(local, remote string) (net.Conn, error) {
		conn, err := dial(local, remote)
		if err != nil {
			return nil, err
		}
		stop = context.AfterFunc(ctx, func() {
			_ = conn.Close()
			close(done)
		})
		return conn, nil
	}
	defer func() {
		if stop != nil && !stop() {
			<-done
		}
	}()

	response, err := query(server, options)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return response, err
}
