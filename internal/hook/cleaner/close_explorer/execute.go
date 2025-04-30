package close_explorer

import (
	"github.com/anyshake/observer/pkg/logger"
)

func (p *CloseExplorerCleanerImpl) Execute() error {
	logger.GetLogger(p.GetName()).Info("closing connection to hardware")
	defer logger.GetLogger(p.GetName()).Info("hardware connection has been closed")
	return p.HardwareDev.Close()
}
