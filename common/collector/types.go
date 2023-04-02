package main

import "com.geophone.observer/common/geophone"

type DataBuffer struct {
	UUID         string
	Station      string
	Latitude     float64
	Longitude    float64
	Altitude     float64
	Syncronized  bool
	Acceleration [10]geophone.Acceleration
}

type SysStatus struct {
	Messages int64
	Errors   int64
	Loss     int64
}
