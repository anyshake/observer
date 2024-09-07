package seisevent

type MagnitudeType string

type Magnitude struct {
	Type  MagnitudeType `json:"type"`
	Value float64       `json:"value"`
}

type Estimation struct {
	P_Wave float64 `json:"p_wave"`
	S_Wave float64 `json:"s_wave"`
}

type Event struct {
	Verfied    bool        `json:"verfied"`
	Timestamp  int64       `json:"timestamp"`
	Event      string      `json:"event"`
	Region     string      `json:"region"`
	Depth      float64     `json:"depth"`
	Latitude   float64     `json:"latitude"`
	Longitude  float64     `json:"longitude"`
	Distance   float64     `json:"distance"`
	Magnitude  []Magnitude `json:"magnitude"`
	Estimation Estimation  `json:"estimation"`
}

type DataSourceProperty struct {
	ID      string            `json:"id"`
	Country string            `json:"country"` // ISO 3166-1 alpha-2 country code
	Deafult string            `json:"default"` // default language key
	Locales map[string]string `json:"locales"` // key: language, value: name
}

type DataSource interface {
	GetProperty() DataSourceProperty
	GetEvents(latitude, longitude float64) ([]Event, error)
}
