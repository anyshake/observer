package seedlink

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/anyshake/observer/feature"
)

type TIME struct{}

// Callback of "TIME <...>" command, implements SeedLinkCommandCallback interface
func (t *TIME) Callback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, streamer SeedLinkStreamer, conn net.Conn, args ...string) error {
	resCode := RES_OK
	switch len(args) {
	case 2:
		endTime, err := t.getTimeFromArg(args[1])
		if err != nil {
			resCode = RES_ERR
		} else {
			cl.EndTime = endTime
		}
		fallthrough
	case 1:
		startTime, err := t.getTimeFromArg(args[0])
		if err != nil {
			resCode = RES_ERR
		} else {
			cl.StartTime = startTime
		}
	default:
		resCode = RES_ERR
	}

	_, err := conn.Write([]byte(resCode))
	return err
}

// Fallback of "TIME <...>" command, implements SeedLinkCommandCallback interface
func (*TIME) Fallback(sl *SeedLinkGlobal, cl *SeedLinkClient, options *feature.FeatureOptions, conn net.Conn, args ...string) {
	conn.Close()
}

func (*TIME) getTimeFromArg(timeStr string) (time.Time, error) {
	if len(timeStr) != 19 {
		return time.Time{}, nil
	}
	splitTimeStr := strings.Split(timeStr, ",")
	if len(splitTimeStr) != 6 {
		return time.Time{}, nil
	}

	// Format:  YYYY,MM,DD,hh,mm,ss
	// Example: 2024,01,16,07,15,16
	year, err := strconv.Atoi(splitTimeStr[0])
	if err != nil {
		return time.Time{}, err
	}

	monthInt, err := strconv.Atoi(splitTimeStr[1])
	if err != nil {
		return time.Time{}, err
	}

	month := time.Month(monthInt)
	day, err := strconv.Atoi(splitTimeStr[2])
	if err != nil {
		return time.Time{}, err
	}

	hour, err := strconv.Atoi(splitTimeStr[3])
	if err != nil {
		return time.Time{}, err
	}

	minute, err := strconv.Atoi(splitTimeStr[4])
	if err != nil {
		return time.Time{}, err
	}

	second, err := strconv.Atoi(splitTimeStr[5])
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(year, month, day, hour, minute, second, 0, time.UTC), nil
}
