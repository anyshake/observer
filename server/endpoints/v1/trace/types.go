package trace

import "time"

type Trace struct{}

const EXPIRATION = time.Minute // Cache expiration duration for calling external API response

type request struct {
	Source string `form:"source" json:"source" xml:"source" binding:"required"`
}
