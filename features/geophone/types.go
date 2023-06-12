package geophone

const (
	PACKET_SIZE   int     = 375
	FILTER_TAPS   int     = 11
	FILTER_CUTOFF float64 = 10.0
)

type Geophone struct {
	Vertical   [PACKET_SIZE]float64
	EastWest   [PACKET_SIZE]float64
	NorthSouth [PACKET_SIZE]float64
	Checksum   [3]uint8
}

type Acceleration struct {
	Timestamp  int64                `json:"timestamp"`
	Vertical   [PACKET_SIZE]float64 `json:"vertical"`
	EastWest   [PACKET_SIZE]float64 `json:"east_west"`
	NorthSouth [PACKET_SIZE]float64 `json:"north_south"`
	Synthesis  [PACKET_SIZE]float64 `json:"synthesis"`
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
