package trace

import "time"

const EXPIRATION = time.Minute // Cache expiration duration for calling external API response

type Trace struct{}

type traceBinding struct {
	Source string `form:"source" json:"source" xml:"source" binding:"required"`
}

type seismicEventEstimation struct {
	P float64 `json:"p"`
	S float64 `json:"s"`
}

type seismicEvent struct {
	Verfied    bool                   `json:"verfied"`
	Timestamp  int64                  `json:"timestamp"`
	Event      string                 `json:"event"`
	Region     string                 `json:"region"`
	Depth      float64                `json:"depth"`
	Latitude   float64                `json:"latitude"`
	Longitude  float64                `json:"longitude"`
	Distance   float64                `json:"distance"`
	Magnitude  float64                `json:"magnitude"`
	Estimation seismicEventEstimation `json:"estimation"`
}

type dataSource interface {
	Property() string
	Fetch() ([]byte, error)
	Parse([]byte) (map[string]any, error)
	List(latitude, longitude float64) ([]seismicEvent, error)
	Format(float64, float64, map[string]any) ([]seismicEvent, error)
}

type dataSourceCache struct {
	Time  time.Time
	Cache []byte
}
