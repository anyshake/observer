package history

type History struct{}

type Binding struct {
	Timestamp int64 `form:"timestamp" json:"timestamp" xml:"timestamp" binding:"required,numeric"`
}
