package history

type History struct{}

type Binding struct {
	Start int64 `form:"start" json:"start" xml:"start" binding:"required,numeric"`
	End   int64 `form:"end" json:"end" xml:"end" binding:"required,numeric"`
}
