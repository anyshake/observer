package publisher

import (
	"time"
)

func Subscribe(gp *Geophone, enable *bool, onMessage func(gp *Geophone) error) {
	lastTime := time.Now().UTC().UnixMilli()

	for *enable {
		if gp.TS != lastTime && len(gp.EHZ) > 0 && len(gp.EHE) > 0 && len(gp.EHN) > 0 {
			err := onMessage(gp)
			if err != nil {
				return
			}

			lastTime = gp.TS
		}

		time.Sleep(10 * time.Millisecond)
	}
}
