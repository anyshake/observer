package geophone

type Geophone struct {
	Vertical [100]int32 `struct:"littleendian"`
}

type Acceleration struct {
	Vertical  [100]float64
	Timestamp int64
}

type ReaderOptions struct {
	Geophone        *Geophone
	Acceleration    *Acceleration
	Sensitivity     float64
	OnErrorCallback func(err error)
	OnDataCallback  func(acceleration *Acceleration)
}
