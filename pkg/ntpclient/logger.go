package ntpclient

// Logger receives synchronization progress and recoverable source failures.
// Implementations must support concurrent calls, return promptly, and not call
// back into the client. No logging is performed unless a logger is provided.
type Logger interface {
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
}

type Option func(*Client)

// WithLogger enables progress logging without coupling the client to a logging library.
func WithLogger(logger Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
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
