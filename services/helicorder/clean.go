package helicorder

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (m *HelicorderService) handleClean() error {
	if m.lifeCycle == 0 {
		return nil
	}

	entries, err := os.ReadDir(m.basePath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".svg") &&
			strings.ContainsAny(entry.Name(), m.stationCode) &&
			strings.ContainsAny(entry.Name(), m.networkCode) {

			info, err := entry.Info()
			if err != nil {
				return err
			}

			duration := time.Duration(m.lifeCycle) * time.Hour * 24
			if time.Now().After(info.ModTime().Add(duration)) {
				err = os.Remove(filepath.Join(m.basePath, entry.Name()))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
