package main

import (
	"fmt"
	"os"
	"path"

	"github.com/anyshake/observer/pkg/logger"
)

func setupLogger(level, logPath string, maxSize, rotation, lifeCycle int) (string, error) {
	if logPath != "" {
		logPath = path.Clean(logPath)

		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			err = os.MkdirAll(logPath, os.ModePerm)
			if err != nil {
				return "", fmt.Errorf("failed to create log directory: %w", err)
			}
		}

		logger.SetFile(logPath, maxSize, rotation, lifeCycle)
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
		return "", fmt.Errorf("unknown log level: %s", level)
	}
	if err != nil {
		return "", fmt.Errorf("failed to set log level: %w", err)
	}

	return logPath, nil
}
