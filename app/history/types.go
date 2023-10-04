package history

import "time"

const (
	UNDEFINED = "-12345  "
	NETWORK   = "OBSERVER"
	DURATION  = time.Hour   // The maximum duration of the data to be queried
	THRESHOLD = time.Minute // There are uneven gaps between the data if time difference is greater than THRESHOLD
)

type History struct{}

type Binding struct {
	Start   int64  `form:"start" json:"start" xml:"start" binding:"required,numeric"`
	End     int64  `form:"end" json:"end" xml:"end" binding:"required,numeric"`
	Format  string `form:"format" json:"format" xml:"format" binding:"required,oneof=json sac"`
	Channel string `form:"channel" json:"channel" xml:"channel" binding:"omitempty,oneof=EHZ EHE EHN"`
}
