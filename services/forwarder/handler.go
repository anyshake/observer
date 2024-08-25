package forwarder

import (
	"fmt"
	"net"
	"strings"

	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/utils/logger"
)

func (a *ForwarderService) handleConnection(conn net.Conn) {
	defer a.unsubscribe(conn.RemoteAddr().String())
	defer conn.Close()

	logger.GetLogger(a.GetServiceName()).Infof("accepted connection from %s", conn.RemoteAddr().String())
	defer logger.GetLogger(a.GetServiceName()).Infof("closed connection from %s", conn.RemoteAddr().String())

	a.subscribe(conn.RemoteAddr().String(), func(data *explorer.ExplorerData) {
		dataBytes := []byte(fmt.Sprintf("$%s,%s,%s,%s,%d,%d,%s,*%02X\r\n$%s,%s,%s,%s,%d,%d,%s,*%02X\r\n$%s,%s,%s,%s,%d,%d,%s,*%02X\r\n",
			a.networkCode, a.stationCode, a.locationCode, fmt.Sprintf("%sZ", a.channelPrefix), data.Timestamp, data.SampleRate, strings.Trim(strings.Replace(fmt.Sprint(data.Z_Axis), " ", ",", -1), "[]"), a.getChecksum(data.Z_Axis),
			a.networkCode, a.stationCode, a.locationCode, fmt.Sprintf("%sE", a.channelPrefix), data.Timestamp, data.SampleRate, strings.Trim(strings.Replace(fmt.Sprint(data.E_Axis), " ", ",", -1), "[]"), a.getChecksum(data.E_Axis),
			a.networkCode, a.stationCode, a.locationCode, fmt.Sprintf("%sN", a.channelPrefix), data.Timestamp, data.SampleRate, strings.Trim(strings.Replace(fmt.Sprint(data.N_Axis), " ", ",", -1), "[]"), a.getChecksum(data.N_Axis),
		))
		_, err := conn.Write(dataBytes)
		if err != nil {
			logger.GetLogger(a.GetServiceName()).Errorln(err)
			return
		}
	})

	for {
		_, err := conn.Read(make([]byte, 1))
		if err != nil {
			return
		}
	}
}
