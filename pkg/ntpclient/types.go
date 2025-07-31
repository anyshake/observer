package ntpclient

import "time"

type Client struct {
	ntpAddr     string
	retries     int
	readTimeout time.Duration
}
