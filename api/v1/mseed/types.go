package mseed

type MSeed struct{}

type mseedBinding struct {
	Action string `form:"action" json:"action" xml:"action" binding:"required,oneof=export show"`
	Name   string `form:"name" json:"name" xml:"name" binding:"omitempty,endswith=.mseed"`
}

type miniSeedFileInfo struct {
	TTL  int    `json:"ttl"`
	Time int64  `json:"time"`
	Size string `json:"size"`
	Name string `json:"name"`
}
