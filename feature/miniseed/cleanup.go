package miniseed

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (m *MiniSEED) handleCleanup(basePath, station, network string, lifeCycle int) {
	for {
		expiredFiles := []string{}
		walkFn := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				modTime := info.ModTime()
				duration := time.Duration(lifeCycle) * time.Hour * 24
				if time.Now().After(modTime.Add(duration)) &&
					strings.HasSuffix(path, ".mseed") &&
					strings.ContainsAny(path, station) &&
					strings.ContainsAny(path, network) {
					expiredFiles = append(expiredFiles, path)
				}
			}

			return nil
		}

		err := filepath.Walk(basePath, walkFn)
		if err != nil {
			m.OnError(nil, err)
		}

		for _, file := range expiredFiles {
			err := os.Remove(file)
			if err != nil {
				m.OnError(nil, err)
			}
		}

		time.Sleep(time.Minute)
	}
}
