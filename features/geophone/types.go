package geophone

type Geophone struct {
	Latitude   float32
	Longitude  float32
	Altitude   float32
	Vertical   float32
	EastWest   float32
	NorthSouth float32
}

type Acceleration struct {
	Timestamp  int64   `json:"timestamp"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Altitude   float64 `json:"altitude"`
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
	Sensitivity     struct {
		Vertical   float64
		EastWest   float64
		NorthSouth float64
	}
	LocationFallback struct {
		Latitude  float64
		Longitude float64
		Altitude  float64
	}
}
