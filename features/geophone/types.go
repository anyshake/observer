package geophone

const (
	FRAME_SIZE int = 100
)

type Geophone struct {
	Vertical   [FRAME_SIZE]float32
	EastWest   [FRAME_SIZE]float32
	NorthSouth [FRAME_SIZE]float32
}

type Acceleration struct {
	Timestamp  int64               `json:"timestamp"`
	Vertical   [FRAME_SIZE]float64 `json:"vertical"`
	EastWest   [FRAME_SIZE]float64 `json:"east_west"`
	NorthSouth [FRAME_SIZE]float64 `json:"north_south"`
	Synthesis  [FRAME_SIZE]float64 `json:"synthesis"`
}

type GeophoneOptions struct {
	Geophone        *Geophone
	Acceleration    *Acceleration
	OnErrorCallback func(error)
	OnDataCallback  func(*Acceleration)
	Latitude        float64
	Longitude       float64
	Altitude        float64
	Sensitivity     struct {
		Vertical   float64
		EastWest   float64
		NorthSouth float64
	}
}
