package ntpclient

import "errors"

var ErrClosed = errors.New("NTP client is closed")

// Close cancels outstanding queries and waits for their probe workers to exit.
// It is safe to call concurrently or more than once. A closed client cannot be reused.
func (c *Client) Close() error {
	if c.cancel == nil {
		return nil
	}
	c.cancel()
	c.gate <- struct{}{}
	<-c.gate
	return nil
}
