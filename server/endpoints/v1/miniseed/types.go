package miniseed

type MiniSEED struct{}

type request struct {
	Action string `form:"action" json:"action" xml:"action" binding:"required,oneof=export list"`
	Name   string `form:"name" json:"name" xml:"name" binding:"omitempty,endswith=.mseed"`
}

type miniSeedFileInfo struct {
	TTL  int    `json:"ttl"`
	Time int64  `json:"time"`
	Size int64  `json:"size"`
	Name string `json:"name"`
}
