package ntpclient

import (
	"time"

	"github.com/bclswl0827/observer/feature"
)

func (n *NTPClient) Start(options *feature.FeatureOptions) {
	var (
		host     = options.Config.NTPClient.Host
		port     = options.Config.NTPClient.Port
		timeout  = options.Config.NTPClient.Timeout
		interval = options.Config.NTPClient.Interval
	)

	n.OnStart(options, "service has started")
	for {
		result, err := n.Read(host, port, timeout)
		if err != nil {
			n.OnError(options, err)
			time.Sleep(time.Second)
			continue
		}

		n.OnReady(options, result)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
