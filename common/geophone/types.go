package geophone

type Geophone struct {
	Latitude  float32
	Longitude float32
	Altitude  float32
	Vertical  [100]int32
}

type Acceleration struct {
	Timestamp int64        `json:"timestamp"`
	Latitude  float64      `json:"latitude"`
	Longitude float64      `json:"longitude"`
	Altitude  float64      `json:"altitude"`
	Vertical  [100]float64 `json:"vertical"`
}

type GeophoneOptions struct {
	Sensitivity      float64
	Geophone         *Geophone
	Acceleration     *Acceleration
	OnErrorCallback  func(error)
	OnDataCallback   func(*Acceleration)
	LocationFallback struct {
		Latitude  float64
		Longitude float64
		Altitude  float64
	}
}
