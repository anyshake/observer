package publisher

type Geophone struct {
	TS  int64   `json:"ts"`
	EHZ []int32 `json:"ehz"`
	EHE []int32 `json:"ehe"`
	EHN []int32 `json:"ehn"`
}

type System struct {
	Messages int64   `json:"messages"`
	Pushed   int64   `json:"pushed"`
	Errors   int64   `json:"errors"`
	Failures int64   `json:"failures"`
	Queued   int64   `json:"queued"`
	Offset   float64 `json:"offset"`
}

type Status struct {
	IsReady  bool      // If false, stuck to wait for time syncing
	Buffer   *Geophone // Buffer area, should not be externally accessed
	System   *System
	Geophone Geophone
}
