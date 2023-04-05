package handler

import (
	"com.geophone.observer/common/geophone"
)

func HandleMessages(options *HandlerOptions, v interface{}) {
	HandleStatus(options, v)

	switch v := v.(type) {
	case *geophone.Acceleration:
		options.Message.Acceleration[options.Status.Messages%10] = *v
	}

	if options.Status.Messages%10 == 0 {
		HandlePush(options)
	}
}
