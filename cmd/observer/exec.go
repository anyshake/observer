//go:build !windows
// +build !windows

package main

import (
	"os"
	"syscall"

	"github.com/anyshake/observer/pkg/logger"
)

func executeBinary(exePath string) {
	err := syscall.Exec(exePath, os.Args, os.Environ())
	if err != nil {
		logger.GetLogger(main).Errorf("process exec failed: %v", err)
	}
}
