package logger

import (
	"errors"

	"github.com/sirupsen/logrus"
)

func SetLevel(level LogLevel) error {
	switch level {
	case INFO:
		logrus.SetLevel(logrus.InfoLevel)
		return nil
	case WARN:
		logrus.SetLevel(logrus.WarnLevel)
		return nil
	case ERROR:
		logrus.SetLevel(logrus.ErrorLevel)
		return nil
	case FATAL:
		logrus.SetLevel(logrus.FatalLevel)
		return nil
	}

	return errors.New("unknown log level")
}
