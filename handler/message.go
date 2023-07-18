package handler

import (
	"time"
)

func OnMessage(gp *Geophone, onMessage func(gp *Geophone) error) {
	lastTime := time.Now().UTC().UnixMilli()

	for {
		var (
			ehz = gp.EHZ
			ehe = gp.EHE
			ehn = gp.EHN
		)
		if gp.TS != lastTime && len(ehz) > 0 && len(ehe) > 0 && len(ehn) > 0 {
			err := onMessage(gp)
			if err != nil {
				return
			}

			lastTime = gp.TS
		}

		time.Sleep(50 * time.Millisecond)
	}
}
