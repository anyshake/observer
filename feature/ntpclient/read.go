package ntpclient

import (
	"time"

	"github.com/beevik/ntp"
)

func (n *NTPClient) read(server string, port, timeout int) (float64, error) {
	response, err := ntp.QueryWithOptions(server, ntp.QueryOptions{
		Port: port, Timeout: time.Duration(time.Duration(timeout).Seconds()),
	})
	if err != nil {
		return 0, err
	}

	return response.ClockOffset.Seconds(), nil
}
