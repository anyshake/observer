package seedlink

import (
	"fmt"
	"net"
	"time"

	"github.com/anyshake/observer/feature"
	"github.com/bclswl0827/mseedio"
	"github.com/clbanning/anyxml"
)

type INFO struct{}

// Callback of "INFO <...>" command, implements SeedLinkCommandCallback interface
func (i *INFO) Callback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, streamer SeedLinkStreamer, conn net.Conn, args ...string) error {
	err := fmt.Errorf("arg error")
	if len(args) < 1 {
		return err
	}

	var (
		action    = args[0]
		dataBytes []byte
	)
	switch action {
	case "ID":
		state := sl.SeedLinkState
		dataBytes, err = i.getID(state, FLAG_INF)
	case "STATIONS":
		var (
			state    = sl.SeedLinkState
			stations = sl.Stations
		)
		dataBytes, err = i.getStations(state, stations)
	case "CAPABILITIES", "CONNECTIONS":
		var (
			state        = sl.SeedLinkState
			capabilities = sl.Capabilities
		)
		dataBytes, err = i.getCapabilities(state, capabilities)
	case "STREAMS":
		var (
			streams  = sl.Streams
			stations = sl.Stations
		)
		dataBytes, err = i.getStreams(sl.SeedLinkState, streams, stations)
	default:
		state := sl.SeedLinkState
		dataBytes, err = i.getID(state, FLAG_ERR)
	}
	if err != nil {
		return err
	}

	_, err = conn.Write(dataBytes)
	return err
}

// Fallback of "INFO <...>" command, implements SeedLinkCommandCallback interface
func (i *INFO) Fallback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, conn net.Conn, args ...string) {
	conn.Write([]byte(RES_ERR))
}

// getID returns response of "INFO ID" command
func (i *INFO) getID(state SeedLinkState, flag int) ([]byte, error) {
	result := map[string]string{
		"-software":     state.Software,
		"-started":      state.StartTime,
		"-organization": state.Organization,
	}
	xmlData, err := anyxml.Xml(result, "seedlink")
	if err != nil {
		return []byte(RES_ERR), err
	}
	// Set XML header and return response
	xmlBody := i.setXMLHeader(xmlData)
	currentTime := time.Now().UTC()
	return i.setResponse(xmlBody, flag, currentTime)
}

// getStations returns response of "INFO STATIONS" command
func (i *INFO) getStations(state SeedLinkState, staions []SeedLinkStation) ([]byte, error) {
	result := map[string]any{
		"-software":     state.Software,
		"-started":      state.StartTime,
		"-organization": state.Organization,
		"station":       staions,
	}
	xmlData, err := anyxml.Xml(result, "seedlink")
	if err != nil {
		return []byte(RES_ERR), err
	}
	// Set XML header and return response
	xmlBody := i.setXMLHeader(xmlData)
	currentTime := time.Now().UTC()
	return i.setResponse(xmlBody, FLAG_INF, currentTime)
}

// getCapabilities returns response of "INFO CAPABILITIES" command
func (i *INFO) getCapabilities(state SeedLinkState, capabilities []SeedLinkCapability) ([]byte, error) {
	result := map[string]any{
		"-software":     state.Software,
		"-started":      state.StartTime,
		"-organization": state.Organization,
		"capability":    capabilities,
	}
	xmlData, err := anyxml.Xml(result, "seedlink")
	if err != nil {
		return []byte(RES_ERR), err
	}
	// Set XML header and return response
	xmlBody := i.setXMLHeader(xmlData)
	currentTime := time.Now().UTC()
	return i.setResponse(xmlBody, FLAG_INF, currentTime)
}

// getStreams returns response of "INFO STREAMS" command
func (i *INFO) getStreams(state SeedLinkState, streams []SeedLinkStream, stations []SeedLinkStation) ([]byte, error) {
	type respModel struct {
		SeedLinkStation
		Streams     []SeedLinkStream `xml:"stream"`
		StreamCheck string           `xml:"stream_check,attr"`
	}
	result := map[any]any{
		"-software":     state.Software,
		"-started":      state.StartTime,
		"-organization": state.Organization,
	}
	var resp []respModel
	for _, v := range stations {
		// Match stream by station name
		var availableStreams []SeedLinkStream
		for _, s := range streams {
			if s.Station == v.Station {
				availableStreams = append(availableStreams, s)
			}
		}
		resp = append(resp, respModel{
			SeedLinkStation: v,
			Streams:         availableStreams,
			StreamCheck:     "enabled",
		})
	}
	result["station"] = resp
	xmlData, err := anyxml.Xml(result, "seedlink")
	if err != nil {
		return []byte(RES_ERR), err
	}
	// Set XML header and return response
	xmlBody := i.setXMLHeader(xmlData)
	currentTime := time.Now().UTC()
	return i.setResponse(xmlBody, FLAG_INF, currentTime)
}

// setXMLHeader sets XML header to body and return string
func (i *INFO) setXMLHeader(body []byte) []byte {
	header := []byte(`<?xml version="1.0" encoding="utf-8"?>`)
	return append(header, body...)
}

// setResponse assembles response in MiniSeed format
func (i *INFO) setResponse(body []byte, errFlag int, startTime time.Time) ([]byte, error) {
	// Convert body to int32 array
	bodyBuffer := []int32{}
	for _, v := range body {
		bodyBuffer = append(bodyBuffer, int32(v))
	}
	// Set channel code by error flag
	channelCode := "INF"
	if errFlag == FLAG_ERR {
		channelCode = "ERR"
	}
	// Initialize MiniSeed data
	var miniseed mseedio.MiniSeedData
	miniseed.Init(mseedio.ASCII, mseedio.MSBFIRST)
	// Split data into 512 bytes each
	bodyLength := len(bodyBuffer)
	dataLength := (512 - mseedio.FIXED_SECTION_LENGTH - mseedio.BLOCKETTE100X_SECTION_LENGTH)
	fullLength := bodyLength + mseedio.FIXED_SECTION_LENGTH + mseedio.BLOCKETTE100X_SECTION_LENGTH
	blockCount := fullLength / 512
	// "SLINFO<space>*" or "SLINFO<space><space>" is signature
	// * indicates non-final block, <space> indicates final block
	blockHeader := []byte{'S', 'L', 'I', 'N', 'F', 'O', ' ', '*'}
	// Append each block to MiniSeed data
	var resultBuffer []byte
	for i := 0; i <= blockCount; i++ {
		startIndex := i * dataLength
		endIndex := (i + 1) * dataLength
		if i == blockCount {
			// Set final block flag
			blockHeader[7] = ' '
			endIndex = bodyLength
		}
		err := miniseed.Append(
			bodyBuffer[startIndex:endIndex],
			&mseedio.AppendOptions{
				SequenceNumber: fmt.Sprintf("%06d", i+1),
				ChannelCode:    channelCode,
				StartTime:      startTime,
				StationCode:    "INFO ",
				LocationCode:   "  ",
				NetworkCode:    "SL",
				SampleRate:     0,
			},
		)
		if err != nil {
			return nil, err
		}
		// Encode MiniSeed data
		res, err := miniseed.Encode(mseedio.APPEND, mseedio.MSBFIRST)
		if err != nil {
			return nil, err
		}
		// Each block should be 512 bytes
		if len(res) < 512 {
			// Fill with 0x00 if length is less than 512
			res = append(res, make([]byte, 512-len(res))...)
		}
		resultBuffer = append(resultBuffer, append(blockHeader, res...)...)
	}
	return resultBuffer, nil
}
