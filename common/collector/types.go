package collector

import "com.geophone.observer/common/geophone"

type Status struct {
	Messages int64
	Errors   int64
	Offset   float64
}

type Message struct {
	UUID         string
	Station      string
	Acceleration [10]geophone.Acceleration
}

type Options struct {
	Status          *Status
	OnDataCallback  func(*Status)
	OnErrorCallback func(error)
}
