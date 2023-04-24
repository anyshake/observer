package ntpclient

import (
	"time"

	"github.com/beevik/ntp"
)

func ReadNTP(server string, options NTPOptions) error {
	response, err := ntp.QueryWithOptions(server, ntp.QueryOptions{
		Port:    options.Port,
		Timeout: time.Duration(time.Duration(options.Timeout).Seconds()),
	})
	if err != nil {
		options.OnErrorCallback(err)
		return err
	}

	options.NTP = &NTP{
		Offset: response.ClockOffset.Seconds(),
	}
	options.OnDataCallback(options.NTP)
	return nil
}

func ReaderDaemon(server string, interval int, options NTPOptions) {
	for {
		ReadNTP(server, options)
		time.Sleep(time.Second * time.Duration(interval))
	}
}
