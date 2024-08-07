package logger

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func SetFile(path string) {
	logrus.AddHook(lfshook.NewHook(
		path, &logrus.JSONFormatter{},
	))
}
