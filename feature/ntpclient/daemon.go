package ntpclient

import (
	"time"

	"com.geophone.observer/feature"
)

func (n *NTPClient) Start(options *feature.FeatureOptions) {
	var (
		host     = options.Config.NTPClient.Host
		port     = options.Config.NTPClient.Port
		timeout  = options.Config.NTPClient.Timeout
		interval = options.Config.NTPClient.Interval
	)

	options.OnStart(MODULE, options, "service has started")
	for {
		result, err := n.Read(host, port, timeout)
		if err != nil {
			options.OnError(MODULE, options, err)
			time.Sleep(time.Second)
			continue
		}

		options.OnReady(MODULE, options, result)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
