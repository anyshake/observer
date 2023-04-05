package geophone

type Geophone struct {
	Latitude  float32
	Longitude float32
	Altitude  float32
	Vertical  [100]int32
}

type Acceleration struct {
	Timestamp int64
	Latitude  float64
	Longitude float64
	Altitude  float64
	Vertical  [100]float64
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
