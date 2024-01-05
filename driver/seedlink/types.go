package seedlink

import (
	"encoding/xml"
	"net"
	"time"

	"github.com/anyshake/observer/feature"
	"github.com/anyshake/observer/publisher"
)

const (
	ORGANIZATION string = "anyshake.org"
	RELEASE      string = "SeedLink v3.1 AnyShake Edition (Very basic implementation in Go) :: SLPROTO:3.1 CAP EXTREPLY NSWILDCARD BATCH WS:13 :: Constructing Realtime Seismic Network Ambitiously."
)

// SeedLink error flags
const (
	FLAG_INF = iota
	FLAG_ERR
)

// SeedLink response codes
const (
	RES_OK  = "OK\r\n"
	RES_ERR = "ERROR\r\n"
)

// SeedLink main daemon config & state
type SeedLinkGlobal struct {
	SeedLinkState
	Streams      []SeedLinkStream
	Stations     []SeedLinkStation
	Capabilities []SeedLinkCapability
}

// SeedLink basic state
type SeedLinkState struct {
	Software     string
	StartTime    string
	Organization string
}

// Station field model of INFO STATIONS command
type SeedLinkStation struct {
	XMLName       xml.Name `xml:"station"`
	BeginSequence string   `xml:"begin_seq,attr"`
	EndSequence   string   `xml:"end_seq,attr"`
	Station       string   `xml:"name,attr"`
	Network       string   `xml:"network,attr"`
	Description   string   `xml:"description,attr"`
}

// Stream field model of INFO STREAMS command
type SeedLinkStream struct {
	XMLName   xml.Name `xml:"stream"`
	BeginTime string   `xml:"begin_time,attr"`
	EndTime   string   `xml:"end_time,attr"`
	SeedName  string   `xml:"seedname,attr"`
	Location  string   `xml:"location,attr"`
	Type      string   `xml:"type,attr"`
	// Exclusive attribute to match station
	Station string `xml:"station,attr"`
}

// Capability field model of INFO CAPABILITY command
type SeedLinkCapability struct {
	XMLName xml.Name `xml:"capability"`
	Name    string   `xml:"name,attr"`
}

// Built-in commands of SeedLink
type SeedLinkCommand struct {
	HasExtraArgs bool
	SeedLinkCommandCallback
}

// SeedLink client state
type SeedLinkClient struct {
	StreamMode bool
	Network    string
	Station    string
	Location   string
	Channels   []string
	StartTime  time.Time
	EndTime    time.Time
}

type SeedLinkStreamer func(pub *publisher.Geophone, conn net.Conn, channels []string, network, station, location string, seqNum *int64) error

// Interface for SeedLink command callback & fallback
type SeedLinkCommandCallback interface {
	Callback(*SeedLinkGlobal, *SeedLinkClient, *feature.FeatureOptions, SeedLinkStreamer, net.Conn, ...string) error
	Fallback(*SeedLinkGlobal, *SeedLinkClient, *feature.FeatureOptions, net.Conn, ...string)
}
