package explorer

import "github.com/anyshake/observer/utils/logger"

type explorerLoggerImpl struct {
	moduleName string
}

func (e *explorerLoggerImpl) Infof(format string, args ...any) {
	logger.GetLogger(e.moduleName).Infof(format, args...)
}

func (e *explorerLoggerImpl) Warnf(format string, args ...any) {
	logger.GetLogger(e.moduleName).Warnf(format, args...)
}

func (e *explorerLoggerImpl) Errorf(format string, args ...any) {
	logger.GetLogger(e.moduleName).Errorf(format, args...)
}
