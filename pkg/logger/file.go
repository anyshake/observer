package logger

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func SetFile(filePath string, maxSize, rotation, lifeCycle int) {
	rotateFileHook := &rotateHook{
		formatter: &logrus.JSONFormatter{},
		logWriter: &lumberjack.Logger{
			Filename:   filePath,
			MaxSize:    maxSize,
			MaxBackups: rotation,
			MaxAge:     lifeCycle,
			Compress:   true,
		},
	}

	logrus.AddHook(rotateFileHook)
}
