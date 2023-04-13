package handlers

import (
	"com.geophone.observer/features/geophone"
	"com.geophone.observer/features/ntpclient"
)

func HandleMessages(options *HandlerOptions, v any) {
	switch v := v.(type) {
	case *ntpclient.NTP:
		options.Status.Offset = v.Offset
	case *geophone.Acceleration:
		v.Timestamp = ntpclient.AlignTime(options.Status.Offset)
		options.Message.Acceleration[options.Status.Messages%10] = *v
		if options.Status.Messages%10 == 0 {
			options.OnReadyCallback(options.Message)
		}
	}

	options.Status.Messages++
}