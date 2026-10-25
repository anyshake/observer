package frp_client

import (
	"context"
	"net"
	"sync"

	"github.com/anyshake/observer/pkg/logger"
	"github.com/fatedier/frp/client"
	v1 "github.com/fatedier/frp/pkg/config/v1"
)

type loggedConnector struct {
	client.Connector
	ctx         context.Context
	controlOnce sync.Once
}

func newLoggedConnector(ctx context.Context, cfg *v1.ClientCommonConfig) client.Connector {
	return &loggedConnector{Connector: client.NewConnector(ctx, cfg), ctx: ctx}
}

func (c *loggedConnector) Connect() (net.Conn, error) {
	conn, err := c.Connector.Connect()
	if err != nil {
		return nil, err
	}

	// FRP opens the control connection first, followed by pooled work connections.
	// Only observe the control connection: work connections routinely reach EOF.
	c.controlOnce.Do(func() {
		conn = &loggedControlConn{Conn: conn, ctx: c.ctx}
	})
	return conn, nil
}

type loggedControlConn struct {
	net.Conn
	ctx       context.Context
	errorOnce sync.Once
}

func (c *loggedControlConn) Read(p []byte) (int, error) {
	n, err := c.Conn.Read(p)
	c.logError("read", err)
	return n, err
}

func (c *loggedControlConn) Write(p []byte) (int, error) {
	n, err := c.Conn.Write(p)
	c.logError("write", err)
	return n, err
}

func (c *loggedControlConn) logError(operation string, err error) {
	if err == nil || c.ctx.Err() != nil {
		return
	}
	c.errorOnce.Do(func() {
		logger.GetLogger(ID).Warnf("FRP control connection %s failed (remote=%s): %v", operation, c.RemoteAddr(), err)
	})
}
