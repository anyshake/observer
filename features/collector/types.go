package collector

import "com.geophone.observer/features/geophone"

type Status struct {
	Messages int64   `json:"messages"`
	Pushed   int64   `json:"pushed"`
	Errors   int64   `json:"errors"`
	Failures int64   `json:"failures"`
	Queued   int64   `json:"queued"`
	Offset   float64 `json:"offset"`
}

type Message struct {
	UUID         string                    `json:"uuid"`
	Station      string                    `json:"station"`
	Latitude     float64                   `json:"latitude"`
	Longitude    float64                   `json:"longitude"`
	Altitude     float64                   `json:"altitude"`
	Acceleration [10]geophone.Acceleration `json:"acceleration"`
}

type CollectorOptions struct {
	Enable             bool
	Status             *Status
	Message            *Message
	OnCompleteCallback func(any)
	OnErrorCallback    func(error)
}
