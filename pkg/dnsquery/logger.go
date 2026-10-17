package dnsquery

// Logger implementations must be concurrent-safe, return promptly, and not call
// back into the client. Logging is disabled unless supplied with WithLogger.
type Logger interface {
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
}

type Option func(*Client)

func WithLogger(logger Logger) Option {
	return func(c *Client) { c.logger = logger }
}

func (c *Client) infof(format string, args ...any) {
	if c.logger != nil {
		c.logger.Infof(format, args...)
	}
}

func (c *Client) warnf(format string, args ...any) {
	if c.logger != nil {
		c.logger.Warnf(format, args...)
	}
}
