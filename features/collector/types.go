package collector

import "com.geophone.observer/features/geophone"

type Status struct {
	Messages int64
	Errors   int64
	Fails    int64
	Queued   int64
	Offset   float64
}

type Message struct {
	UUID         string                    `json:"uuid"`
	Station      string                    `json:"station"`
	Acceleration [10]geophone.Acceleration `json:"acceleration"`
}

type CollectorOptions struct {
	Enable             bool
	Status             *Status
	Message            *Message
	OnCompleteCallback func(interface{})
	OnErrorCallback    func(error)
}
