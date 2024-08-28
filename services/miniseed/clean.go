package miniseed

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (m *MiniSeedService) handleClean() error {
	if m.lifeCycle == 0 {
		return nil
	}

	expiredFiles := []string{}
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			modTime := info.ModTime()
			duration := time.Duration(m.lifeCycle) * time.Hour * 24
			if time.Now().After(modTime.Add(duration)) &&
				strings.HasSuffix(path, ".mseed") &&
				strings.ContainsAny(path, m.stationCode) &&
				strings.ContainsAny(path, m.networkCode) {
				expiredFiles = append(expiredFiles, path)
			}
		}

		return nil
	}

	err := filepath.Walk(m.basePath, walkFn)
	if err != nil {
		return err
	}

	for _, file := range expiredFiles {
		err := os.Remove(file)
		if err != nil {
			return err
		}
	}

	return nil
}
