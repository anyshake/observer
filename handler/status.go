package handler

import (
	"com.geophone.observer/common/geophone"
	"com.geophone.observer/common/ntpclient"
)

func HandleStatus(options *HandlerOptions, v interface{}) {
	switch v := v.(type) {
	case *ntpclient.NTP:
		options.Status.Offset = v.Offset
	case *geophone.Acceleration:
		HandleTimestamp(options, v)
	}

	options.Status.Messages++
}
