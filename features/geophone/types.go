package geophone

type Geophone struct {
	Vertical   float32
	EastWest   float32
	NorthSouth float32
}

type Acceleration struct {
	Timestamp  int64   `json:"timestamp"`
	Vertical   float64 `json:"vertical"`
	EastWest   float64 `json:"east_west"`
	NorthSouth float64 `json:"north_south"`
	Synthesis  float64 `json:"synthesis"`
}

type GeophoneOptions struct {
	Geophone        *Geophone
	Acceleration    *Acceleration
	OnErrorCallback func(error)
	OnDataCallback  func(*Acceleration)
	Interval        int
	Latitude        float64
	Longitude       float64
	Altitude        float64
	Sensitivity     struct {
		Vertical   float64
		EastWest   float64
		NorthSouth float64
	}
}
