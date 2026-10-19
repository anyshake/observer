package close_explorer

import (
	"context"

	"github.com/anyshake/observer/pkg/logger"
)

func (p *CloseExplorerCleanerImpl) Execute(context.Context) error {
	if p.HardwareDev != nil {
		logger.GetLogger(p.GetName()).Infoln("closing connection to hardware")
		defer logger.GetLogger(p.GetName()).Infoln("hardware connection has been closed")
		return p.HardwareDev.Close()
	}

	return nil
}
