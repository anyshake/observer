package handler

import (
	"time"

	"com.geophone.observer/common/geophone"
)

func HandleTimestamp(options *HandlerOptions, acceleration *geophone.Acceleration) {
	current := time.Now().UnixMilli()
	offset := options.Status.Offset * float64(time.Second.Milliseconds())

	acceleration.Timestamp = current + int64(offset)
}
