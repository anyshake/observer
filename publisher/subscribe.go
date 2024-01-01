package publisher

import (
	"time"
)

func Subscribe(gp *Geophone, exp *bool, onMessage func(gp *Geophone) error) {
	lastTime := time.Now().UTC().UnixMilli()

	for *exp {
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

		time.Sleep(10 * time.Millisecond)
	}
}
