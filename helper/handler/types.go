package handler

import (
	"com.geophone.observer/common/collector"
	"com.geophone.observer/common/geophone"
)

type HandlerOptions struct {
	Error           error
	Status          *collector.Status
	Message         *collector.Message
	Acceleration    *geophone.Acceleration
	OnReadyCallback func(*collector.Message)
}
