package transport

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

func New(endpoint string, timeout int) (ITransport, error) {
	urlObj, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse transport endpoint: %w", err)
	}

	if timeout == 0 {
		timeout = 5 // 5 seconds timeout by default
	}

	switch urlObj.Scheme {
	case "serial":
		urlObj, err := url.Parse(endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to parse serial endpoint: %w", err)
		}
		deviceName := urlObj.Hostname()
		if deviceName == "" {
			deviceName = urlObj.Path
		}
		baudrate, err := strconv.Atoi(urlObj.Query().Get("baudrate"))
		if err != nil {
			return nil, fmt.Errorf("failed to parse serial baudrate: %w", err)
		}
		return &SerialTransportImpl{baudrate: baudrate, port: deviceName, timeout: time.Duration(timeout) * time.Second}, nil
	case "tcp":
		return &TcpTransportImpl{host: urlObj.Host, timeout: time.Duration(timeout) * time.Second}, nil
	}

	return nil, fmt.Errorf("transport scheme %s is not supported", urlObj.Scheme)
}
