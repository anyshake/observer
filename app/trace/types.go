package trace

import "time"

const (
	CAPACITY   = 30          // Capacity for rate limiter to prevent from being attacked
	EXPIRATION = time.Minute // Cache expiration duration for calling external API response
)

type Trace struct{}

type Binding struct {
	Source string `form:"source" json:"source" xml:"source" binding:"required"`
}

type Event struct {
	Verfied   bool    `json:"verfied"`
	Timestamp int64   `json:"timestamp"`
	Event     string  `json:"event"`
	Region    string  `json:"region"`
	Depth     float64 `json:"depth"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Distance  float64 `json:"distance"`
	Magnitude float64 `json:"magnitude"`
	Estimated float64 `json:"estimated"`
}

type DataSource interface {
	Fetch() ([]byte, error)
	Property() (string, string)
	Parse([]byte) (map[string]any, error)
	List(latitude, longitude float64) ([]Event, error)
	Format(float64, float64, map[string]any) ([]Event, error)
}

type DataSourceCache struct {
	Time  time.Time
	Cache []byte
}
