package explorer

import (
	"github.com/anyshake/observer/cleaners"
	"github.com/anyshake/observer/drivers/explorer"
	"github.com/anyshake/observer/utils/logger"
)

func (p *ExplorerCleanerTask) Execute(options *cleaners.Options) {
	options.Dependency.Invoke(func(explorerDeps *explorer.ExplorerDependency) {
		logger.GetLogger(p.GetTaskName()).Info("closing connection to hardware device")
		explorerDeps.Transport.Close()
	})
}
