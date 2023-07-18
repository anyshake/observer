package handler

import (
	"time"

	"com.geophone.observer/config"
)

func InitHandler(config *config.Conf, status *Status) {
	status.Geophone = Geophone{}
	status.System = &System{}
	status.Buffer = &Geophone{
		EHZ: []int32{},
		EHE: []int32{},
		EHN: []int32{},
		TS:  time.Now().UTC().UnixMilli(),
	}
}
