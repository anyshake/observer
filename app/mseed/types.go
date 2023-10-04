package mseed

type MSeed struct{}

type Binding struct {
	Action string `form:"action" json:"action" xml:"action" binding:"required,oneof=export show"`
	Name   string `form:"name" json:"name" xml:"name" binding:"omitempty,endswith=.mseed"`
}

type MiniSEEDFile struct {
	Time string `json:"time"`
	Size string `json:"size"`
	Name string `json:"name"`
}
