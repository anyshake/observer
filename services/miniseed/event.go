package miniseed

import (
	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/utils/logger"
)

func (m *MiniSeedService) handleExplorerEvent(data *explorer.ExplorerData) {
	m.miniseedBuffer[len(m.miniseedBuffer)-m.writeBufferCountDown] = *data

	m.writeBufferCountDown--
	m.cleanUpCountDown--

	if m.writeBufferCountDown == 0 {
		err := m.handleWrite()
		if err != nil {
			logger.GetLogger(m.GetServiceName()).Warnln(err)
		} else {
			logger.GetLogger(m.GetServiceName()).Infof("%d record(s) has been written to MiniSEED file", m.writeBufferInterval)
		}
		m.writeBufferCountDown = m.writeBufferInterval
	}

	if m.cleanUpCountDown == 0 {
		err := m.handleClean()
		if err != nil {
			logger.GetLogger(m.GetServiceName()).Warnln(err)
		}
		m.cleanUpCountDown = MINISEED_CLEANUP_INTERVAL
	}
}
