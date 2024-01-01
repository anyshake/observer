package publisher

type Int32Array []int32

type Geophone struct {
	TS  int64      `json:"ts" gorm:"ts;index;not null"`
	EHZ Int32Array `json:"ehz" gorm:"ehz;type:text;not null"`
	EHE Int32Array `json:"ehe" gorm:"ehe;type:text;not null"`
	EHN Int32Array `json:"ehn" gorm:"ehn;type:text;not null"`
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
