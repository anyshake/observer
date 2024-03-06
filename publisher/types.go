package publisher

import "time"

type Int32Array []int32

type Geophone struct {
	TS  int64      `json:"ts" gorm:"ts;index;not null"`
	EHZ Int32Array `json:"ehz" gorm:"ehz;type:text;not null"`
	EHE Int32Array `json:"ehe" gorm:"ehe;type:text;not null"`
	EHN Int32Array `json:"ehn" gorm:"ehn;type:text;not null"`
}

type System struct {
	Messages int64   `json:"messages"`
	Errors   int64   `json:"errors"`
	Offset   float64 `json:"offset"`
}

type Status struct {
	LastRecvTime time.Time // Timestamp of last received packet
	ReadyTime    time.Time // If is zero, app will stuck to wait for time syncing
	Geophone     Geophone  // Geophone data of nearest 1 second
	Buffer       *Geophone // Buffer area, should not be externally accessed
	System       *System
}

type ChannelSegmentBuffer struct {
	DataBuffer []int32
	Samples    int32
	SeqNum     int64
}

type SegmentBuffer struct {
	TimeStamp     time.Time
	ChannelBuffer map[string]*ChannelSegmentBuffer
}
