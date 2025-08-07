package logger

import (
	"io"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type fileHookImpl struct {
	formatter logrus.Formatter
	logWriter io.Writer
}

func (hook *fileHookImpl) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *fileHookImpl) Fire(entry *logrus.Entry) (err error) {
	b, err := hook.formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = hook.logWriter.Write(b)
	return err
}

func RegisterFileLogger(filePath string, maxSize, rotation, lifeCycle int) {
	fileHook := &fileHookImpl{
		formatter: &logrus.JSONFormatter{},
		logWriter: &lumberjack.Logger{
			Filename:   filePath,
			MaxSize:    maxSize,
			MaxBackups: rotation,
			MaxAge:     lifeCycle,
			Compress:   true,
		},
	}

	logrus.AddHook(fileHook)
}
