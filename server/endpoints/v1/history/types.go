package history

import "time"

type History struct{}

const (
	JSON_MAX_DURATION = time.Hour   // The maximum duration of the JSON data to be exported
	SAC_MAX_DURATION  = time.Hour   // The maximum duration of the SAC data to be exported
	THRESHOLD         = time.Minute // There are uneven gaps between the data if time difference is greater than THRESHOLD
)

type request struct {
	StartTime int64  `form:"start_time" json:"start_time" xml:"start_time" binding:"required,numeric"`
	EndTime   int64  `form:"end_time" json:"end_time" xml:"end_time" binding:"required,numeric"`
	Format    string `form:"format" json:"format" xml:"format" binding:"required,oneof=json sac miniseed"`
	Channel   string `form:"channel" json:"channel" xml:"channel"`
}
