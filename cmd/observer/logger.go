package main

import (
	"fmt"
	"os"
	"path"

	"github.com/anyshake/observer/pkg/logger"
	"github.com/anyshake/observer/pkg/ringbuf"
)

func setupLogger(level, logPath string, maxSize, rotation, lifeCycle int) (*ringbuf.Buffer[string], string, error) {
	if logPath != "" {
		logPath = path.Clean(logPath)

		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			err = os.MkdirAll(logPath, os.ModePerm)
			if err != nil {
				return nil, "", fmt.Errorf("failed to create log directory: %w", err)
			}
		}

		logger.RegisterFileLogger(logPath, maxSize, rotation, lifeCycle)
	}

	var err error
	switch level {
	case "info":
		err = logger.SetLevel(logger.INFO)
	case "warn":
		err = logger.SetLevel(logger.WARN)
	case "error":
		err = logger.SetLevel(logger.ERROR)
	default:
		return nil, "", fmt.Errorf("unknown log level: %s", level)
	}
	if err != nil {
		return nil, "", fmt.Errorf("failed to set log level: %w", err)
	}

	logBuffer := logger.RegisterBufferLogger(256)
	return logBuffer, logPath, nil
}
