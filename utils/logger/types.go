package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

type LogLevel int

const (
	INFO LogLevel = iota
	WARN
	ERROR
	FATAL
)

const TIMESTAMP_FORMAT = "2006-01-02 15:04:05.000"

type rotateHook struct {
	formatter logrus.Formatter
	logWriter io.Writer
}
