package dnsquery

import "errors"

var ErrClosed = errors.New("DNS client is closed")

// Close permanently cancels the client, joins active probes, and releases idle
// connections. It is safe to call concurrently or more than once.
func (c *Client) Close() error {
	c.cancel()
	c.gate <- struct{}{}
	defer func() { <-c.gate }()
	var errs []error
	for _, state := range c.servers {
		errs = append(errs, state.server.Close())
	}
	return errors.Join(errs...)
}
