package handler

import (
	"com.geophone.observer/features/collector"
	"com.geophone.observer/features/geophone"
)

type HandlerOptions struct {
	Error           error
	Status          *collector.Status
	Message         *collector.Message
	Acceleration    *geophone.Acceleration
	OnReadyCallback func(*collector.Message)
}
