package trace

import "time"

const EXPIRATION = time.Minute // Cache expiration duration for calling external API response

type Trace struct{}

type traceBinding struct {
	Source string `form:"source" json:"source" xml:"source" binding:"required"`
}
