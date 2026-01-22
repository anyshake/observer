package main

import (
	"os"
	"os/exec"

	"github.com/anyshake/observer/pkg/logger"
)

func executeBinary(exePath string) {
	cmd := exec.Command(exePath, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		logger.GetLogger(main).Errorf("spawn self failed: %v", err)
	}
}
