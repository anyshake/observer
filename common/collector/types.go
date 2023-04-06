package collector

import "com.geophone.observer/common/geophone"

type Status struct {
	Messages int64
	Errors   int64
	Queued   int64
	Offset   float64
}

type Message struct {
	UUID         string                    `json:"uuid"`
	Station      string                    `json:"station"`
	Acceleration [10]geophone.Acceleration `json:"acceleration"`
}

type CollectorOptions struct {
	Status             *Status
	Message            *Message
	OnCompleteCallback func(interface{})
	OnErrorCallback    func(error)
}
