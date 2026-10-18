package seedlink

import (
	"net"
	"time"
)

type seedLinkConn struct {
	net.Conn
}

func (c *seedLinkConn) Write(data []byte) (int, error) {
	// Start the timeout at the actual write, after history queries and idle
	// periods. OnData runs after Write and cannot revive an expired deadline.
	if err := c.Conn.SetWriteDeadline(time.Now().Add(seedLinkWriteTimeout)); err != nil {
		return 0, err
	}
	return c.Conn.Write(data)
}
