package handler

import (
	"com.geophone.observer/common/geophone"
	"com.geophone.observer/common/ntpclient"
)

func HandleMessages(options *HandlerOptions, v interface{}) {
	switch v := v.(type) {
	case *geophone.Acceleration:
		HandleTimestamp(options, v)
		options.Message.Acceleration[options.Status.Messages%10] = *v
		if options.Status.Messages%10 == 0 {
			options.OnReadyCallback(options.Message)
		}
	case *ntpclient.NTP:
		options.Status.Offset = v.Offset
	}

	options.Status.Messages++
}
