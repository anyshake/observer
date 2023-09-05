package publisher

import (
	"time"

	"github.com/bclswl0827/observer/config"
)

func Initialize(config *config.Conf, status *Status) {
	status.System = &System{}
	status.Buffer = &Geophone{
		EHZ: []int32{},
		EHE: []int32{},
		EHN: []int32{},
		TS:  time.Now().UTC().UnixMilli(),
	}
	status.Geophone = Geophone{}
}
